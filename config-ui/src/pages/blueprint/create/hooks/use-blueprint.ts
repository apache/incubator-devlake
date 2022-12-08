import { useContext } from 'react'

import { BlueprintContext } from './blueprint-context'

export const useBlueprint = () => {
  return useContext(BlueprintContext)
}
