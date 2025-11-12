import { axiosInstance } from "../../connection/http/axiosInstance"

export async function uploadFileToR2(file) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await axiosInstance.post('/api/messages/upload-image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return response.data
}

export async function getImagefromKey(key) {
  const response = await axiosInstance.get(`api/messages/${key}`)

  return response.data
}