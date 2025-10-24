import { httpService } from '@/services/http.service';

export const UserService = {
  async list(): Promise<User[]> {
    const { data } = await httpService('user/', 'get');
    return data;
  },

  async get(id: number): Promise<User> {
    const { data } = await httpService(`user/${id}`, 'get');
    return data;
  },

  async delete(id: string): Promise<void> {
    await httpService(`user/${id}`, 'delete');
  },

  async create(user: User): Promise<User> {
    const { data } = await httpService('user/', 'post', user);
    return data;
  },

  async update(user: User): Promise<User> {
    const { data } = await httpService(`user/${user.id}`, 'put', user);
    return data;
  },

  async save(user: User): Promise<User> {
    return user.id ? UserService.update(user) : UserService.create(user);
  },
};
