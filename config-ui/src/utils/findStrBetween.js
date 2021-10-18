export const findStrBetween = (str, first, sliceAmount) => {
  const r = new RegExp(first + '(.*?)' + ';', 'gm')

  if (str) {
    return str.match(r)[0].slice(sliceAmount, -1).split(',')
  }
}
