import { setEnvValue } from '../../../../utils/envValue'

export default function handler(req, res) {
  const { key, value } = req.query

  setEnvValue(key, value)

  res.status(200).json({ key: value, status: 'updated' })
}
