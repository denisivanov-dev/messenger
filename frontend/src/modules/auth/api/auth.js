export async function registerUser(payload) {
    const res = await fetch('http://localhost:8000/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
  
    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.detail || 'Ошибка регистрации')
    }
  
    return await res.json()
}