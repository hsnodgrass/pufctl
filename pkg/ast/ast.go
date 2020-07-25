package ast

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

// PuppetfileLexer is a custom regex lexer for Puppetfiles
var PuppetfileLexer = lexer.Must(regex.New(`
	Keyword = ^mod
	String = '([^']*)'
	Int = \d
	Char = [[:alpha:]]
	Ident = :[A-Za-z0-9_]+
	Comment = #.*
	Punct = [,@:.]
	Assign = (=>)
	Whitespace = [\s\t]+
`))

// Parse parses a Puppetfile and returns the parsed AST
func Parse(text string) (*Puppetfile, error) {
	parser, err := NewParser()
	if err != nil {
		return nil, err
	}
	puppetfile := &Puppetfile{}
	err = parser.ParseString(text, puppetfile)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return nil, err
	}
	_ = puppetfile.ParseMetadata()
	_ = puppetfile.SortByName()
	return puppetfile, nil
}

// NewParser creates and returns a Puppetfile parser
func NewParser() (*participle.Parser, error) {
	return participle.Build(
		&Puppetfile{},
		participle.Lexer(PuppetfileLexer),
		participle.Elide("Whitespace"),
		participle.Unquote("String"),
		participle.UseLookahead(3),
	)
}

// Value Struct that holds the types that a value can be
type Value struct {
	Pos    lexer.Position
	String string `@String`
	Ident  string `| @Ident`
}

// Sprint returns string representation of type
func (v *Value) Sprint() string {
	switch {
	case v.String != "":
		return fmt.Sprintf("'%s'", v.String)
	case v.Ident != "":
		return fmt.Sprintf("%s", v.Ident)
	}
	return ""
}

// Property Struct that hold key-value properties pairs
type Property struct {
	Pos   lexer.Position
	Key   *Value `@@ `
	Value *Value `( Assign @@ )?(",")?`
}

// Sprint returns string representation of type
func (p *Property) Sprint() string {
	if p.Value != nil {
		return fmt.Sprintf("%s => %s,", p.Key.Sprint(), p.Value.Sprint())
	}
	return fmt.Sprintf("%s", p.Key.Sprint())
}

// Module Struct that holds each Module and related properties
type Module struct {
	Pos        lexer.Position
	Name       string      `Keyword @String","?`
	Properties []*Property `( @@ )*`
}

// Sprint returns string representation of type
func (m *Module) Sprint() string {
	var modString string
	var propStrings []string
	for i := 0; i < len(m.Properties); i++ {
		propStrings = append(
			propStrings,
			fmt.Sprintf("  %s\n", m.Properties[i].Sprint()),
		)
	}
	fmtPropStrings := strings.Join(propStrings, "")
	if fmtPropStrings != "" {
		modString = fmt.Sprintf("%s '%s',\n", "mod", m.Name)
	} else {
		modString = fmt.Sprintf("%s '%s'\n", "mod", m.Name)
	}

	return fmt.Sprintf("%s%s", modString, fmtPropStrings)
}

// Comment holds the text and lexer position of a comment in the AST
type Comment struct {
	Pos  lexer.Position
	Text string `@Comment`
}

// MetaTag returns the metadata tag in the comment text, if present
func (c *Comment) MetaTag() string {
	re := regexp.MustCompile(`#\s?@(?P<Tag>[\w_-]+):.*`)
	match := re.FindStringSubmatch(c.Text)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// MetaData returns the metadata tag in the comment text, if present
func (c *Comment) MetaData() string {
	re := regexp.MustCompile(`#\s?@[\w_-]+:\s*(?P<Data>.*)`)
	match := re.FindStringSubmatch(c.Text)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// Statement holds statements that make up the Puppetfile, either modules or comments
type Statement struct {
	Pos     lexer.Position
	Module  *Module  `@@`
	Comment *Comment `| @@`
}

// Sprint returns string representation of type
func (s *Statement) Sprint() string {
	if s.Module != nil {
		return s.Module.Sprint()
	} else if s.Comment != nil {
		return s.Comment.Text
	} else {
		return "\n"
	}
}

// Puppetfile Struct that holds entire Puppetfile
type Puppetfile struct {
	Statements     []*Statement `{ @@ }`
	Metadata       Metadata
	ModuleMetadata []ModuleMetadata
}

// ParseMetadata goes through all statements in the Puppetfile
// and makes metadata to module associations. It then populates
// the module's Metadata field.
func (p *Puppetfile) ParseMetadata() error {
	meta := Metadata{MetaPairs: make([]MetaPair, 0)}
	for _, s := range p.Statements {
		if s.Comment != nil {
			tag := s.Comment.MetaTag()
			data := s.Comment.MetaData()
			if tag != "" {
				pair := MetaPair{Tag: tag, Data: data}
				meta.MetaPairs = append(meta.MetaPairs, pair)
			}
		}
	}
	p.Metadata = meta
	return nil
}

// SortByName reorders all Statements to be sorted alphabetically by name.
// Comment statements are reordered based on their original
// positions relative to the modules below them in the Puppetfile.
func (p *Puppetfile) SortByName() error {
	var modules []*Module
	// This gets wierd because we only want to sort the modules
	// by name. So we build an array of only the modules, run
	// ByName sort on that array, and then replace each module
	// in the Statements array with it's sorted counterpart
	modComments := map[*Module][]*Statement{}
	var cmts []*Statement
	var topBlockCmts []*Statement
	lastCmtLine := 0
	for _, s := range p.Statements {
		if s.Comment != nil {
			if s.Comment.Pos.Line == 1 || s.Comment.Pos.Line == lastCmtLine+1 {
				topBlockCmts = append(topBlockCmts, s)
				lastCmtLine = s.Comment.Pos.Line
			} else {
				cmts = append(cmts, s)
			}
		}
		if s.Module != nil {
			if len(cmts) > 0 {
				modComments[s.Module] = append(modComments[s.Module], cmts...)
				meta := Metadata{MetaPairs: make([]MetaPair, 0)}
				for _, c := range cmts {
					tag := c.Comment.MetaTag()
					data := c.Comment.MetaData()
					pair := MetaPair{Tag: tag, Data: data}
					meta.MetaPairs = append(meta.MetaPairs, pair)
				}
				modmeta := ModuleMetadata{Name: s.Module.Name, Metadata: meta}
				p.ModuleMetadata = append(p.ModuleMetadata, modmeta)
				cmts = cmts[:0]
			}
			modules = append(modules, s.Module)
		}
	}
	sort.Sort(ByName(modules))
	// Remove all the comments
	var stmts []*Statement
	stmts = append(stmts, p.Statements...)
	stmts = p.RemoveComments(stmts)
	modIdx := 0
	commentIdx := map[int]*Module{}
	for i := 0; i < len(stmts); i++ {
		if stmts[i].Module != nil {
			if len(modComments[modules[modIdx]]) > 0 {
				commentIdx[i] = modules[modIdx]
			}
			stmts[i].Module = modules[modIdx]
			modIdx++
		}
	}
	offset := 0
	cIdxKeys := make([]int, 0)
	for k := range commentIdx {
		cIdxKeys = append(cIdxKeys, k)
	}
	sort.Ints(cIdxKeys)
	for _, k := range cIdxKeys {
		mc := modComments[commentIdx[k]]
		for _, s := range mc {
			stmts = p.AddComment(k+offset, s, stmts)
			offset++
		}
	}
	p.Statements = p.Statements[:0]
	for _, s := range stmts {
		p.Statements = append(p.Statements, s)
	}
	// Add the top block comments to the top of the Puppetfile
	if len(topBlockCmts) > 0 {
		// This is a hack to get a newline after the top block comments
		dummyPos := lexer.Position{Filename: "dummy", Offset: -1, Line: -1, Column: -1}
		dummyCmt := &Comment{dummyPos, ""}
		dummyStmt := &Statement{dummyPos, nil, dummyCmt}
		topBlockCmts = append(topBlockCmts, dummyStmt)
		p.Statements = append(topBlockCmts, p.Statements...)
	}
	// Add in the leftover comments to the bottom of the Puppetfile
	if len(cmts) > 0 {
		p.Statements = append(p.Statements, cmts...)
	}
	return nil
}

// RemoveComments removes all Statements from the Puppetfile that have a Comment
func (p *Puppetfile) RemoveComments(stmts []*Statement) []*Statement {
	comments := false
	idx := 0
	for i, s := range stmts {
		if s.Comment != nil {
			comments = true
			idx = i
		}
	}
	if comments {
		return p.RemoveComments(append(stmts[:idx], stmts[idx+1:]...))
	}
	return stmts
}

// AddComment adds a Statement with a Comment above the given index
func (p *Puppetfile) AddComment(idx int, s *Statement, stmts []*Statement) []*Statement {
	if idx != 0 {
		return append(stmts[:idx], append([]*Statement{s}, stmts[idx:]...)...)
	}
	return append([]*Statement{s}, stmts...)
}

// AddModuleMetadata adds a new metadata statement above a module in the Puppetfile
func (p *Puppetfile) AddModuleMetadata(name string, tag string, data string) error {
	for i, s := range p.Statements {
		if s.Module != nil && s.Module.Name == name {
			dummyPos := lexer.Position{Filename: "dummy", Offset: -1, Line: -1, Column: -1}
			cmt := &Comment{dummyPos, fmt.Sprintf("# @%s: %s", tag, data)}
			meta := &Statement{dummyPos, nil, cmt}
			stmts := p.AddComment(i, meta, p.Statements)
			p.Statements = p.Statements[:0]
			p.Statements = append(p.Statements, stmts...)
		}
	}
	return nil
}

// SearchModulesByMetaTag returns a slice of module name strings that
// have the given tag associated with them.
func (p *Puppetfile) SearchModulesByMetaTag(tag string) []string {
	finds := make([]string, 0)
	for _, m := range p.ModuleMetadata {
		mp := m.SearchByTag(tag)
		if len(mp) > 0 {
			finds = append(finds, m.Name)
		}
	}
	return finds
}

// SearchModulesByMetaData returns a slice of module name strings that
// have the given metadata associated with them.
func (p *Puppetfile) SearchModulesByMetaData(data string) []string {
	finds := make([]string, 0)
	for _, m := range p.ModuleMetadata {
		mp := m.SearchByData(data)
		if len(mp) > 0 {
			finds = append(finds, m.Name)
		}
	}
	return finds
}

// SearchMetaByModuleName returns a string of all metadata for the given module name
func (p *Puppetfile) SearchMetaByModuleName(name string) string {
	for _, m := range p.ModuleMetadata {
		if m.Name == name {
			return m.Metadata.Sprint()
		}
	}
	return ""
}

// Sprint returns string representation of type
func (p *Puppetfile) Sprint() string {
	var modStrings []string
	for _, s := range p.Statements {
		modStrings = append(modStrings, s.Sprint())
	}
	return fmt.Sprintf("%s\n", strings.Join(modStrings, "\n"))
}

// ByName implements sort.Interface base on the Module struct's Name field
type ByName []*Module

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
