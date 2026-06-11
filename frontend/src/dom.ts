export function inputValue(id: string) {
    return (document.getElementById(id) as HTMLInputElement).value
}

export function domForm(id: string) {
    return (document.getElementById(id) as HTMLFormElement)
}

export function domDiv(id: string) {
    return (document.getElementById(id) as HTMLDivElement)
}

export function domUl(id: string) {
    return (document.getElementById(id) as HTMLUListElement)
}

