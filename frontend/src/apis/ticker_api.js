import axios from 'axios';

const HOST = `http://${window.location.host}`

export const getTickerItems = () => axios.get(`${HOST}/api/ticker`)
export const addTickerItem= (item) => axios.post(`${HOST}/api/ticker`, item)
export const deleteTickerItem = (id) => axios.delete(`${HOST}/api/ticker/${id}`)
