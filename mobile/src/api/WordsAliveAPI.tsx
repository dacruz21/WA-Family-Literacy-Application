import axios, { AxiosInstance } from 'axios';

import { Book, BookDetails } from '../models/Book';
import { Language } from '../models/Languages';
import { User, UpdateUser } from '../models/User';

// Class to encapsulate the handler for the Words Alive API
class WordsAliveAPI {
  client: AxiosInstance;

  constructor(baseURL: string) {
    this.client = axios.create({ baseURL: baseURL });
  }

  // Set the Firebase token for future API calls
  setToken(token: string): void {
    this.client.defaults.headers.Authorization = `Bearer ${token}`;
  }

  // Unset the Firebase token
  clearToken(): void {
    delete this.client.defaults.headers.Authorization;
  }

  // Pings the backend to wake it up if asleep. Does not throw or return anything
  ping(): void {
    this.client.get('/ping').catch(() => {});
  }

  // makes a call to the database and returns an array of all books
  async getBooks(): Promise<Book[]> {
    const res = await this.client.get('/books');
    return res.data;
  }

  // returns an individual book by id
  async getBook(id: string, lang: Language): Promise<BookDetails> {
    const res = await this.client.get(`/books/${id}/${lang}`);
    return res.data;
  }

  async getUser(id: string): Promise<User> {
    const res = await this.client.get(`/users/${id}`);
    return res.data;
  }

  async createUser(user: User): Promise<User> {
    const res = await this.client.post('/users', user);
    return res.data;
  }

  async updateUser(id: string, update: UpdateUser): Promise<User> {
    const res = await this.client.patch(`/users/${id}`, update);
    return res.data;
  }
}

export { WordsAliveAPI };
