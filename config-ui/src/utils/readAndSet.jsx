import { findStrBetween } from './findStrBetween'

export const readAndSet = (tagName, tagLen, isStatus, str, fn1, fn2) => {
  if (isStatus) {
    const strValuesReq = findStrBetween(str, 'Requirement:', ';')
    const strValuesRes = findStrBetween(str, 'Resolved:', ';')

    if (strValuesReq) fn1(strValuesReq[0].slice(12, -1).split(','))
    if (strValuesRes) fn2(strValuesRes[0].slice(9, -1).split(','))
  } else {
    const strValues = findStrBetween(str, tagName, ';')

    if (strValues) fn1(strValues[0].slice(tagLen, -1).split(','))
  }
}
