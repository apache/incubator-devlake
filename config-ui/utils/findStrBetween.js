export const findStrBetween = (str, first, last) => {
  const r = new RegExp(first + '(.*?)' + last, 'gm')
  return str.match(r)
}
