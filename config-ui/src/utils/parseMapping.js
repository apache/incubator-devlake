export function parseMapping (mappingString) {
  const mapping = {}
  if (!mappingString.trim()) {
    return mapping
  }
  for (const item of mappingString.split(';')) {
    let [standard, customs] = item.split(';')
    standard = standard.trim()
    mapping[standard] = mapping[standard] || []
    if (!customs) {
      continue
    }
    for (const custom of customs.split(',')) {
      mapping[standard].push(custom.trim())
    }
  }
  return mapping
}
