import axios from 'axios';

const HOST = `http://${window.location.host}`

export const getEnabledFeatureSlides = () => axios.get(`${HOST}/api/settings/featureslides`)
export const setEnabledFeatureSlides = (slides) => axios.put(`${HOST}/api/settings/featureslides`, slides)
