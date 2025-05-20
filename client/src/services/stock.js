import api from './api';

export const searchStocks = async (query) => {
  const response = await api.post('/stock/search', { query });
  return response.data;
};

export const getStockDetails = async (symbol) => {
  const response = await api.get(`/stock/details?symbol=${symbol}`);
  return response.data;
};

export const getUserStocks = async () => {
  const response = await api.get('/stock/list');
  return response.data;
};

export const addStock = async (symbol) => {
  const response = await api.post('/stock/add', { symbol });
  return response.data;
}; 