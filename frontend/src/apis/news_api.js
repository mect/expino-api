import axios from 'axios';

const HOST = `http://${window.location.host}`

export const getAllNews = () => axios.get(`${HOST}/api/news`)
export const getNews = (id) => axios.get(`${HOST}/api/news/${id}`)
export const addNews = (info) => axios.post(`${HOST}/api/news`, info)
export const editNews = (info) => axios.put(`${HOST}/api/news`, info)
export const deleteNews = (id) => axios.delete(`${HOST}/api/news/${id}`)
