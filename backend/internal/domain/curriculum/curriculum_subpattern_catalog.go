package curriculum

var goalSubpatternSpecs = mergeGoalSubpatternSpecs(
	goalSubpatternSpecsInstrumentLanguageBody,
	goalSubpatternSpecsCraftMakerCooking,
	goalSubpatternSpecsVisualDigitalWriting,
	goalSubpatternSpecsKnowledge,
)

func mergeGoalSubpatternSpecs(parts ...map[string]goalSubpatternSpec) map[string]goalSubpatternSpec {
	merged := make(map[string]goalSubpatternSpec)
	for _, part := range parts {
		for key, spec := range part {
			merged[key] = spec
		}
	}
	return merged
}
