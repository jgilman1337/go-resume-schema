package go_resume_schema

// DefaultExportFormat is the export shape and encoding a resume without a
// settings block exports in: one document per resume section, JSON encoded.
const DefaultExportFormat = SettingsExportFormatFragmentedJson

// DefaultSectionOrder is the default print order for a resume, the junior
// technical preset: the summary under the header, then the skills list and
// experience in the top-left hotspot a reviewer reads first, with projects
// and education following. A document's settings reorder these same
// sections; every section missing from a stored ordering appends in this
// order. Basics never appears because the header always prints first.
var DefaultSectionOrder = []Section{
	SectionSummary,
	SectionSkills,
	SectionWork,
	SectionProjects,
	SectionEducation,
	SectionCertificates,
	SectionPublications,
	SectionAwards,
	SectionVolunteer,
	SectionLanguages,
	SectionInterests,
	SectionReferences,
}

// NewGradSectionOrder is the new grad preset: education leads because the
// degree is the strongest evidence a recent graduate carries, then the
// skills list and projects, with internships and early work after. A
// document adopting this preset copies its sections verbatim; sections it
// omits still append in the default junior order.
var NewGradSectionOrder = []Section{
	SectionSummary,
	SectionEducation,
	SectionSkills,
	SectionProjects,
	SectionWork,
	SectionCertificates,
	SectionAwards,
	SectionPublications,
	SectionVolunteer,
	SectionLanguages,
	SectionInterests,
	SectionReferences,
}

// SeniorSectionOrder is the senior technical preset: experience and its
// accomplishments lead the technical evidence, the skills list follows for
// keyword coverage, and education precedes the optional projects tail. A
// document adopting this preset copies its sections verbatim; sections it
// omits still append in the default junior order.
var SeniorSectionOrder = []Section{
	SectionSummary,
	SectionWork,
	SectionSkills,
	SectionEducation,
	SectionProjects,
	SectionCertificates,
	SectionPublications,
	SectionAwards,
	SectionVolunteer,
	SectionLanguages,
	SectionInterests,
	SectionReferences,
}

// NonTechnicalSectionOrder is the non-technical preset: experience leads
// right after the summary, the skills list and education follow, and the
// community sections (certificates, volunteer) outrank the technical tail
// (projects, publications). A document adopting this preset copies its
// sections verbatim; sections it omits still append in the default junior
// order.
var NonTechnicalSectionOrder = []Section{
	SectionSummary,
	SectionWork,
	SectionSkills,
	SectionEducation,
	SectionCertificates,
	SectionVolunteer,
	SectionProjects,
	SectionAwards,
	SectionPublications,
	SectionLanguages,
	SectionInterests,
	SectionReferences,
}

// DefaultSettings builds the settings block a resume without one behaves
// with: the junior technical print ordering, no exclusions, and the
// fragmented JSON export format. The caller owns the returned block.
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
