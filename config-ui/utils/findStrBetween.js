export const findStrBetween = (str, first, last) => {
  const r = new RegExp(first + '(.*?)' + last, 'gm')

  if (str) {
    return str.match(r)
  }
}
