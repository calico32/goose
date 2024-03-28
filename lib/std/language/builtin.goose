export const int = frozen {
  parse: native "F/int.parse",
  tryParse: native "F/int.tryParse",
  // min: native "C/int.min",
  // max: native "C/int.max",
}

export const float = frozen {
  parse: native "F/float.parse",
  tryParse: native "F/float.tryParse",
  // max: native "C/float.max",
  // min: native "C/float.min",
  // infinity: native "C/float.infinity",
  // nan: native "C/float.nan",
}

export const string = frozen {
  // format: native "F/string.format",
}

export const bool = frozen {
  parse: native "F/bool.parse",
  // tryParse: native "F/bool.tryParse",
}
