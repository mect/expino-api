import axios from 'axios';

const HOST = `http://${window.location.host}`

export const getGraphItems = () => axios.get(`${HOST}/api/graphs`)
export const addGraphItem= (item) => axios.post(`${HOST}/api/graphs`, item)
export const deleteGraphItem = (id) => axios.delete(`${HOST}/api/graphs/${id}`)
