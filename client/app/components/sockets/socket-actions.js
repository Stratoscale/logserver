import {setEntities} from 'entities/entities-actions';

export function updateBatchEntry(entities, batchModel) {
  return (dispatch, getState) => {
    dispatch(setEntities(entities, batchModel.entity.getKey(), batchModel.entity.getIdAttribute(), {merge: true}));
  }
}