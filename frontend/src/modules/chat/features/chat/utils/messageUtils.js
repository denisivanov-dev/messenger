import { uploadFileToR2 } from "../api/chatApi"

export async function uploadAllImages(files) {
    const MAX_FILES = 5

    if (files.length > MAX_FILES) {
        throw new Error(`Максимум можно загрузить ${MAX_FILES} изображений`)
    }
  
    const results = []

    for (const file of files) {
        const { key, size } = await uploadFileToR2(file)
        results.push({
        key,
        type: 'image',
        size,
        original_name: file.name
        })
    }

    return results
}