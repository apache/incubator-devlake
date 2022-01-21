import * as dayjs from 'dayjs'
import * as relativeTime from 'dayjs/plugin/relativeTime'
import * as updateLocale from 'dayjs/plugin/updateLocale'
import * as LocalizedFormat from 'dayjs/plugin/localizedFormat'
import * as utc from 'dayjs/plugin/utc'

const localeConfiguration = {
  relativeTime: {
    future: 'in %s',
    past: '%s ago',
    s: '< 1min',
    m: 'a minute',
    mm: '%d minutes',
    h: 'an hour',
    hh: '%d hours',
    d: 'a day',
    dd: '%d days',
    M: 'a month',
    MM: '%d months',
    y: 'a year',
    yy: '%d years'
  }
}

dayjs.extend(relativeTime)
dayjs.extend(updateLocale)
dayjs.extend(LocalizedFormat)
dayjs.extend(utc)
dayjs.updateLocale('en', localeConfiguration)

export default dayjs
