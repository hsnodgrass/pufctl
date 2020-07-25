package ast

import (
	"fmt"
	"sort"
	"strings"
)

// Metadata holds an array of MetaPair objects
type Metadata struct {
	MetaPairs []MetaPair
}

// AllTags returns a sorted array of all tag strings
func (m *Metadata) AllTags() []string {
	mp := make([]string, 0)
	for _, m := range m.MetaPairs {
		if m.Tag != "" {
			mp = append(mp, m.Tag)
		}
	}
	sort.Strings(mp)
	return mp
}

// TagExists checks if a tag exists in Metadata and returns the
// an array of indexs of the tag MetaPairs and a boolean if it
// exists or not. The index return value of a non existent tag is -1.
func (m *Metadata) TagExists(tag string) ([]int, bool) {
	indexes := make([]int, 0)
	for i, elem := range m.MetaPairs {
		if elem.Tag == tag {
			indexes = append(indexes, i)
		}
	}
	if len(indexes) > 0 {
		return indexes, true
	}
	return []int{-1}, false
}

// DataExists checks if data exists in Metadata and returns the
// an array of indexs of the data MetaPairs and a boolean if it
// exists or not. The index return value of non existent data is -1.
func (m *Metadata) DataExists(data string) ([]int, bool) {
	indexes := make([]int, 0)
	for i, elem := range m.MetaPairs {
		if elem.Data == data {
			indexes = append(indexes, i)
		}
	}
	if len(indexes) > 0 {
		return indexes, true
	}
	return []int{-1}, false
}

// SearchByTag returns all metapairs with the specified tag or nil
func (m *Metadata) SearchByTag(tag string) []MetaPair {
	finds := make([]MetaPair, 0)
	indexes, exists := m.TagExists(tag)
	if exists {
		for _, i := range indexes {
			finds = append(finds, m.MetaPairs[i])
		}
	}
	if len(finds) > 0 {
		return finds
	}
	return nil
}

// SearchByData returns all metapairs with the specified data or nil
func (m *Metadata) SearchByData(data string) []MetaPair {
	finds := make([]MetaPair, 0)
	indexes, exists := m.DataExists(data)
	if exists {
		for _, i := range indexes {
			finds = append(finds, m.MetaPairs[i])
		}
	}
	if len(finds) > 0 {
		return finds
	}
	return nil
}

// Sprint returns a string representation of the object
func (m *Metadata) Sprint() string {
	metas := make([]string, 0)
	for _, mp := range m.MetaPairs {
		metas = append(metas, mp.Sprint())
	}
	return strings.Join(metas, "\n")
}

// MetaPair holds the actual metadata, all tags with that metadata,
// and all modules tagged with that metadata.
type MetaPair struct {
	Tag  string
	Data string
}

// Sprint returns a string of the MetaPair
func (m *MetaPair) Sprint() string {
	return fmt.Sprintf("# @%s: %s", m.Tag, m.Data)
}

// ModuleMetadata holds the module name and all the metadata assigned to the module.
type ModuleMetadata struct {
	Name     string
	Metadata Metadata
}

// SearchByTag returns all metapairs with the specified tag or nil
func (m *ModuleMetadata) SearchByTag(tag string) []MetaPair {
	return m.Metadata.SearchByTag(tag)
}

// SearchByData returns all metapairs with the specified data or nil
func (m *ModuleMetadata) SearchByData(data string) []MetaPair {
	return m.Metadata.SearchByData(data)
}
