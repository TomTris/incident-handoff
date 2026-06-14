import type { UserContext } from '@/types';

export const makeEmptyUserContext = () : UserContext => {return {id: '', username : '', role : ''}}

// Format of s: 2026-06-14T16:28:53.926576+02:00
export function isoToDateAndTime(s: string) : string {
    const date = s.slice(11, 16) + " " + s.slice(8, 10) + "." + s.slice(5, 7)
    return date
}