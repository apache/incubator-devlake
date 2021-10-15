import { getEnvValue } from '../../../../utils/envValue'

export default function handler(req, res) {
  const { key } = req.query

  res.status(200).json({ key: getEnvValue(key), status: 'done' })
}
