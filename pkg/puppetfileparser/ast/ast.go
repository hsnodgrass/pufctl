// Package ast provides the abstract syntax tree, parser, and everything else in puppetfileparser
package ast

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"

	"github.com/asaskevich/govalidator"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

// PuppetfileLexer is a custom regex lexer for Puppetfiles
var PuppetfileLexer = lexer.Must(regex.New(`
	Keyword = ^mod
	Forge = ^forge
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
	err = puppetfile.ParseMetadata()
	if err != nil {
		return nil, err
	}
	err = puppetfile.SortByName()
	if err != nil {
		return nil, err
	}
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
	default:
		return ""
	}
}

// Checksum returns an md5 checksum of the String or Ident contained in the value.
func (v Value) Checksum() [16]byte {
	var data []byte
	if v.String != "" {
		data = []byte(v.String)
	} else if v.Ident != "" {
		data = []byte(v.Ident)
	} else {
		data = []byte("__null__")
	}
	return md5.Sum(data)
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

// Checksum returns an md5 checksum of the Sprint() output
func (p Property) Checksum() [16]byte {
	data := []byte(p.Sprint())
	return md5.Sum(data)
}

// ChecksumKey returns the md5 checksum of the Key field
func (p Property) ChecksumKey() [16]byte {
	return p.Key.Checksum()
}

// ChecksumValue returns the md5 checksum of the Value field
func (p Property) ChecksumValue() [16]byte {
	return p.Value.Checksum()
}

// OverwriteValue creates a new *Value adds it to the Property
// Value field.
func (p *Property) OverwriteValue(val string) {
	var value *Value
	if strings.HasPrefix(val, ":") {
		value = &Value{
			Ident: val,
		}
	} else {
		value = &Value{
			String: val,
		}
	}
	p.Value = value
}

// Module Struct that holds each Module and related properties
type Module struct {
	Pos        lexer.Position
	Name       string      `Keyword @String","?`
	Properties []*Property `( @@ )*`
}

// GetProperty returns a pointer to a Property of the module,
// if the Module has said property. GetProperty performs this
// search by Property key.
func (m Module) GetProperty(key string) *Property {
	for _, p := range m.Properties {
		if p.Key != nil {
			if p.Key.Ident != "" && p.Key.Ident == key {
				return p
			}
		}
	}
	return nil
}

// EditProperty overwrites the value of an existing property
func (m *Module) EditProperty(key, value string) {
	if key == "bare" {
		m.NewBareProperty(value)
	} else {
		prop := m.GetProperty(key)
		if prop == nil {
			m.AddProperties([]string{fmt.Sprintf("%s=>%s", key, value)})
		} else {
			prop.OverwriteValue(value)
		}
	}
}

// NewBareProperty overwrites all other properties in a module
// with a single bare property.
func (m *Module) NewBareProperty(value string) {
	var val *Value
	if strings.HasPrefix(value, ":") {
		val = &Value{Ident: value}
	} else {
		val = &Value{String: value}
	}
	prop := &Property{Key: val}
	m.Properties = []*Property{prop}
}

// GetPropertyValue accepts a string of a param name and returns the
// associated value if present. As Puppetfiles are distinctly
// idomatic, there are some inputs that can potentially match
// several values:
// Input "version": Returns a bare version string, the :latest
//     symbol as a string, :tag, or :ref if :ref is a semver.
// Input "branch": Returns :branch, :default_branch, or :ref
//     if :ref is a git branch.
func (m Module) GetPropertyValue(name string) string {
	versionKeys := []string{"version", ":tag", ":ref"}
	branchKeys := []string{":branch", ":default_branch", ":ref"}
	props := map[string]string{}
	for _, p := range m.Properties {
		key := strings.Replace(p.Key.Sprint(), "'", "", -1)
		if p.Value != nil {
			val := strings.Replace(p.Value.Sprint(), "'", "", -1)
			props[key] = val
		} else {
			props["version"] = key
		}
	}
	switch name {
	case "version":
		for _, k := range versionKeys {
			if v, found := props[k]; found {
				if k == ":ref" && govalidator.IsSemver(v) {
					return v
				} else if k == ":ref" {
					continue
				}
				return v
			}
		}
	case "branch":
		for _, k := range branchKeys {
			if v, found := props[k]; found {
				if k == ":ref" && govalidator.Matches(k, "[A-Za-z0-9-]+") {
					return v
				} else if k == ":ref" {
					continue
				}
				return v
			}
		}
	}
	if v, found := props[name]; found {
		return v
	}
	return ""
}

// AddProperties accepts strings in the form of "key=>value",
// parses them into Property objects, and adds them to the module.
// Key must have a ":" prefixing it, just like in a Puppetfile.
func (m *Module) AddProperties(props []string) {
	for _, prop := range props {
		p := &Property{DummyPos(), nil, nil}
		match := MapNamedCaptureGroups(ReAssignment, prop)
		if key, ok := match["Key"]; ok {
			keyVal := &Value{Pos: DummyPos(), String: "", Ident: key}
			p.Key = keyVal
		}
		if val, ok := match["Value"]; ok {
			if val != "" {
				valVal := &Value{Pos: DummyPos(), String: val, Ident: ""}
				p.Value = valVal
			}
		}
		if p.Key != nil {
			m.Properties = append(m.Properties, p)
		} else {
			m.AddProperty(props[0])
		}
	}
}

// AddProperty accepts a string and adds it as a property to the module.
// AddProperty is used to add special properties (bare version string / :latest symbol)
func (m *Module) AddProperty(prop string) {
	p := &Property{DummyPos(), nil, nil}
	match := ReIdent.FindString(prop)
	if match != "" {
		keyVal := &Value{Pos: DummyPos(), String: "", Ident: match}
		p.Key = keyVal
	} else {
		keyVal := &Value{Pos: DummyPos(), String: prop, Ident: ""}
		p.Key = keyVal
	}
	if p.Key != nil {
		m.Properties = append(m.Properties, p)
	}
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

// Checksum returns an md5 checksum of the Sprint() output
func (m Module) Checksum() [16]byte {
	data := []byte(m.Sprint())
	return md5.Sum(data)
}

// ChecksumName returns an md5 checksum of the Module Name
func (m Module) ChecksumName() [16]byte {
	data := []byte(m.Name)
	return md5.Sum(data)
}

// ChecksumProperties returns an array of [k, v] arrays where
// k = Property Key field checksum and v= Property Value field checksum.
func (m Module) ChecksumProperties() [][][16]byte {
	propSums := [][][16]byte{}
	for _, p := range m.Properties {
		key := p.ChecksumKey()
		val := p.ChecksumValue()
		propSums = append(propSums, [][16]byte{key, val})
	}
	return propSums
}

// Comment holds the text and lexer position of a comment in the AST
type Comment struct {
	Pos  lexer.Position
	Text string `@Comment`
}

// Checksum returns the md5 checksum of the Text field
func (c Comment) Checksum() [16]byte {
	data := []byte(c.Text)
	return md5.Sum(data)
}

// Tag returns the metadata tag in the comment text, if present
func (c *Comment) Tag() string {
	return TagFromString(c.Text)
}

// Data returns the metadata data in the comment text, if present
func (c *Comment) Data() string {
	return DataFromString(c.Text)
}

// MetaPair returns a MetaPair if the Comment's text contains one
func (c *Comment) MetaPair() (MetaPair, error) {
	mp, err := PairFromString(c.Text)
	if err != nil {
		return MetaPair{}, err
	}
	return mp, nil
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

// Checksum returns an md5 checksum of the Sprint() output
func (s Statement) Checksum() [16]byte {
	data := []byte(s.Sprint())
	return md5.Sum(data)
}

// Forge holds a Forge declaration
type Forge struct {
	URL string `Forge @String`
}

// Sprint returns a string representation of type
func (f *Forge) Sprint() string {
	return fmt.Sprintf("forge '%s'", f.URL)
}

// Checksum returns an md5 checksum of the Sprint() output
func (f Forge) Checksum() [16]byte {
	data := []byte(f.Sprint())
	return md5.Sum(data)
}

// Puppetfile Struct that holds entire Puppetfile
type Puppetfile struct {
	Forge               *Forge       `( @@ )?`
	Statements          []*Statement `{ @@ }`
	Metadata            Metadata
	ModuleMetadata      []ModuleMetadata
	TopBlockComments    []*Statement
	BottomBlockComments []*Statement
	ModuleVersionMap    map[string]string
}

// HasModule returns true if the named module exists in the Puppetfile
func (p Puppetfile) HasModule(name string) (bool, int) {
	stmts := []*Statement{}
	stmts = append(stmts, p.Statements...)
	stmts = p.RemoveComments(stmts)
	sort.Sort(StmtByName(stmts))
	idx := sort.Search(len(stmts), func(i int) bool { return stmts[i].Module.Name >= name })
	if idx < len(stmts) && stmts[idx].Module.Name == name {
		return true, idx
	}
	return false, -1
}

// GetModule returns a Module struct by name
func (p Puppetfile) GetModule(name string) *Module {
	stmts := []*Statement{}
	stmts = append(stmts, p.Statements...)
	stmts = p.RemoveComments(stmts)
	sort.Sort(StmtByName(stmts))
	idx := sort.Search(len(stmts), func(i int) bool { return stmts[i].Module.Name >= name })
	if idx < len(stmts) && stmts[idx].Module.Name == name {
		return stmts[idx].Module
	}
	return nil
}

// RenameModule renames a module in the Puppetfile while keeping
// the module's properties the same. The Puppetfile is then sorted
// by name.
func (p *Puppetfile) RenameModule(name, new string) error {
	mod := p.GetModule(name)
	if mod == nil {
		return fmt.Errorf("Module %s can't be found in the Puppetfile", name)
	}
	mod.Name = new
	p.SortByName()
	return nil
}

// ParseMetadata goes through all statements in the Puppetfile
// and makes metadata to module associations. It then populates
// the module's Metadata field.
func (p *Puppetfile) ParseMetadata() error {
	meta := Metadata{MetaPairs: make([]MetaPair, 0)}
	for _, s := range p.Statements {
		if s.Comment != nil {
			mp, err := s.Comment.MetaPair()
			if err != nil {
				continue
			}
			meta.MetaPairs = append(meta.MetaPairs, mp)
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
	p.ModuleVersionMap = map[string]string{}
	modComments := map[*Module][]*Statement{}
	var cmts []*Statement
	lastCmtLine := 0
	for _, s := range p.Statements {
		if s.Comment != nil {
			if s.Comment.Pos.Line == 1 || s.Comment.Pos.Line == lastCmtLine+1 {
				p.TopBlockComments = append(p.TopBlockComments, s)
				lastCmtLine = s.Comment.Pos.Line
			} else {
				cmts = append(cmts, s)
			}
		}
		if s.Module != nil {
			p.ModuleVersionMap[s.Module.Name] = s.Module.GetPropertyValue("version")
			if len(cmts) > 0 {
				modComments[s.Module] = append(modComments[s.Module], cmts...)
				meta := Metadata{MetaPairs: make([]MetaPair, 0)}
				for _, c := range cmts {
					mp, err := c.Comment.MetaPair()
					if err != nil {
						continue
					}
					meta.MetaPairs = append(meta.MetaPairs, mp)
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
			stmts = p.AddStatementAbove(k+offset, s, stmts)
			offset++
		}
	}
	p.Statements = p.Statements[:0]
	for _, s := range stmts {
		p.Statements = append(p.Statements, s)
	}
	// Add in the leftover comments to the bottom of the Puppetfile
	if len(cmts) > 0 {
		p.BottomBlockComments = append(p.BottomBlockComments, cmts...)
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

// AddModule adds a module and its properties to the Puppetfile
func (p *Puppetfile) AddModule(slug string, props []string) error {
	if f, _ := p.HasModule(slug); !f {
		dMod := DummyModule()
		dMod.Module.Name = slug
		dMod.Module.AddProperties(props)
		err := p.AddStatement(dMod)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Puppetfile already contains module %s", slug)
}

// AddComment adds a Comment to either the top or bottom block comments
func (p *Puppetfile) AddComment(location, text string) error {
	com := DummyComment()
	com.Comment.Text = text
	switch location {
	case "top":
		p.TopBlockComments = append(p.TopBlockComments, com)
	case "bottom":
		p.BottomBlockComments = append(p.BottomBlockComments, com)
	default:
		return fmt.Errorf("Location %s is not valid, should be \"top\" or \"bottom\"", location)
	}
	return nil
}

// AddStatement appends a Statement the Puppetfiles Statements property and sorts Statements
func (p *Puppetfile) AddStatement(s *Statement) error {
	var stmts []*Statement
	stmts = append(stmts, p.Statements...)
	p.Statements = p.Statements[:0]
	stmts = append(stmts, s)
	p.Statements = append(p.Statements, stmts...)
	err := p.ParseMetadata()
	if err != nil {
		return err
	}
	err = p.SortByName()
	if err != nil {
		return err
	}
	return nil
}

// AddStatementAbove adds a Statement to a slice of Statements above the given index
func (p *Puppetfile) AddStatementAbove(idx int, s *Statement, stmts []*Statement) []*Statement {
	if idx != 0 {
		return append(stmts[:idx], append([]*Statement{s}, stmts[idx:]...)...)
	}
	return append([]*Statement{s}, stmts...)
}

// AddModuleMetadata adds a new metadata statement above a module in the Puppetfile
func (p *Puppetfile) AddModuleMetadata(name string, tag string, data string) error {
	for i, s := range p.Statements {
		if s.Module != nil && s.Module.Name == name {
			var stmts []*Statement
			stmts = append(stmts, p.Statements...)
			p.Statements = p.Statements[:0]
			meta := DummyComment()
			meta.Comment.Text = fmt.Sprintf("# @%s: %s", tag, data)
			stmts = p.AddStatementAbove(i, meta, stmts)
			p.Statements = append(p.Statements, stmts...)
			return nil
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
	var topBlock []string
	var modStrings []string
	var bottomBlock []string
	for _, t := range p.TopBlockComments {
		topBlock = append(topBlock, t.Sprint())
	}
	for _, s := range p.Statements {
		modStrings = append(modStrings, s.Sprint())
	}
	for _, b := range p.BottomBlockComments {
		bottomBlock = append(bottomBlock, b.Sprint())
	}
	if p.Forge != nil {
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n", strings.Join(topBlock, "\n"), p.Forge.Sprint(), strings.Join(modStrings, "\n"), strings.Join(bottomBlock, "\n"))
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s\n", strings.Join(topBlock, "\n"), strings.Join(modStrings, "\n"), strings.Join(bottomBlock, "\n"))
}

// Checksum returns an md5 checksum of the stringified contents of Puppetfile
func (p Puppetfile) Checksum() [16]byte {
	data := []byte(p.Sprint())
	return md5.Sum(data)
}

// ByName implements sort.Interface base on the Module struct's Name field
type ByName []*Module

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// StmtByName implements sort.Interface base on slice of Statements base on Module field's Name field
type StmtByName []*Statement

func (a StmtByName) Len() int           { return len(a) }
func (a StmtByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StmtByName) Less(i, j int) bool { return a[i].Module.Name < a[j].Module.Name }

// DummyPos returns a lexer.Position with fake data
func DummyPos() lexer.Position {
	return lexer.Position{Filename: "dummy", Offset: -1, Line: -1, Column: -1}
}

// DummyStatement returns a Statement pointer with nil fields
func DummyStatement() *Statement {
	return &Statement{DummyPos(), nil, nil}
}

// DummyModule returns a Statement pointer with a dummy Module field
func DummyModule() *Statement {
	dummyStmt := DummyStatement()
	dummyStmt.Module = &Module{DummyPos(), "", make([]*Property, 0)}
	return dummyStmt
}

// DummyComment returns a Statement pointer with a dummy Comment field
func DummyComment() *Statement {
	dummyStmt := DummyStatement()
	dummyStmt.Comment = &Comment{DummyPos(), ""}
	return dummyStmt
}
