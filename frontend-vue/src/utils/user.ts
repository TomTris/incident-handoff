import type { UserContext } from '@/types';

export const makeEmptyUserContext = () : UserContext => {return {id: '', username : '', role : ''}}