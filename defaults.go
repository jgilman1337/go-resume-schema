package go_resume_schema

// DefaultExportFormat is the export shape and encoding a resume without a
// settings block exports in: one document per resume section, JSON encoded.
const DefaultExportFormat = SettingsExportFormatFragmentedJson

// DefaultSectionOrder is the default print order for a resume: the summary
// under the header, then skills and education ahead of experience, the
// ordering that suits a technical resume. A document's settings reorder
// these same sections; every section missing from a stored ordering appends
// in this order. Basics never appears because the header always prints
// first.
var DefaultSectionOrder = []Section{
	SectionSummary,
	SectionSkills,
	SectionEducation,
	SectionWork,
	SectionProjects,
	SectionCertificates,
	SectionPublications,
	SectionAwards,
	SectionVolunteer,
	SectionLanguages,
	SectionInterests,
	SectionReferences,
}

// DefaultSettings builds the settings block a resume without one behaves
// with: the technical print ordering, no exclusions, and the fragmented
// JSON export format. The caller owns the returned block.
//
//	@return The default settings block.
func DefaultSettings() *Settings {
	// The ordering copies so callers mutating one block never leak into
	// another.
	order := make([]Section, len(DefaultSectionOrder))
	copy(order, DefaultSectionOrder)

	// The exclusions start empty rather than nil so the block marshals.
	return &Settings{
		SectionOrder: order,
		Exclusions:   []SettingsExclusionsElem{},
		ExportFormat: DefaultExportFormat,
	}
}
