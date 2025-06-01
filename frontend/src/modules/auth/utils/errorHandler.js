export const applyErrors = (errors, fieldsMap) => {
   Object.entries(errors).forEach(([key, message]) => {
      if (fieldsMap[key]) {
         fieldsMap[key].value = message
      }
   })
}