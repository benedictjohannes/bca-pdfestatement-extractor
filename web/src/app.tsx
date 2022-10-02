import { useState, useEffect } from 'preact/hooks'

type textTypes = {
    title: string
    subTitle: string
    dropTitle: string
    dropSubTitle: string
    remarks: string
    errNotPdf: string
    tryAgain: string
}
const textId: textTypes = {
    title: 'Export Excel',
    subTitle: 'untuk BCA E-Statement ',
    dropTitle: 'Klik untuk pilih PDF eStatement',
    dropSubTitle: 'atau drag and drop PDF eStatement ke sini',
    remarks:
        'Untuk menjamin privasi Anda, PDF eStatement Anda diproses langsung dalam halaman ini tanpa upload ke server. Data Anda tidak akan keluar dari komputer Anda untuk export Excel ini.',
    errNotPdf: 'File bukan PDF',
    tryAgain: 'Silakan Coba Lagi',
}
const textEn: textTypes = {
    title: 'Excel Export',
    subTitle: 'for BCA E-Statement ',
    dropTitle: 'Click to pick eStatement PDF',
    dropSubTitle: 'or drag and drop your eStatement PDF here',
    remarks:
        'For your privacy, your eStatement PDF is processed on this page without any uploading to server. Any of your data never leave your browser.',
    errNotPdf: "File isn't PDF",
    tryAgain: 'Try again',
}
const addActive = () =>
    document.getElementById('dropper')?.classList.add('active')
const rmActive = () => {
    document.getElementById('dropper')?.classList.remove('active')
    document.getElementById('dropper')?.classList.remove('dragErr')
}

const onDragEnter = (e: DragEvent) => {
    console.log({ onDragStart: e })
    addActive()
}
const onDragEnd = (e: DragEvent) => {
    rmActive()
}

const onDragOver = (e: Event) => e.preventDefault()

const Loading = () => (
    <div className='loadingWrap'>
        <div className='loadingDiv' />
    </div>
)
const WasmErr = () => (
    <div className='loadingWrap'>
        <h3 className='wasmErr'>Error loading Excel Export WASM library</h3>
    </div>
)

const processFiles = async (files: File[]) => {
    const bufs = await Promise.all(files.map(f => f.arrayBuffer()))
    const fileBytes = bufs.map(r => new Uint8Array(r))
    files.forEach((file, i) => {
        const fileName = file.name
        const fileByte = fileBytes[i]
        // @ts-ignore
        const result = window.excelExport(fileName, fileByte) as string
        console.log({ result })
        if (result.startsWith('ERR:')) return
        const el = document.createElement('a')
        el.style.display = 'none'
        el.download = fileName.replace('.pdf', '.xlsx')
        document.body.appendChild(el)
        const byteChars = atob(result)
        const byteNumbers = new Array(byteChars.length)
        for (let i = 0; i < byteChars.length; i++) {
            byteNumbers[i] = byteChars.charCodeAt(i)
        }
        const blob = new Blob([new Uint8Array(byteNumbers)], {
            type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        })
        const url = URL.createObjectURL(blob)
        el.href = url
        el.click()
    })
}

const wasmUrl = new URL('./exportExcelPDFeStatement.wasm', import.meta.url)

const App = () => {
    const [wasmErr, setWasmErr] = useState(false)
    const [loading, setLoading] = useState(true)
    useEffect(() => {
        // @ts-ignore
        const go = new Go()
        fetch(wasmUrl)
            .then(r => r.arrayBuffer())
            .then(wasmBin => WebAssembly.instantiate(wasmBin, go.importObject))
            .then(wasmLoadResult => {
                go.run(wasmLoadResult.instance)
            })
            .catch(e => {
                console.error({ wasmLoadErr: e })
                setWasmErr(true)
            })
            .finally(() => setLoading(false))
    }, [])
    const [err, setErr] = useState('')
    useEffect(() => {
        if (!!err) {
            document.getElementById('dropper')?.classList.add('error')
        } else {
            document.getElementById('dropper')?.classList.remove('error')
        }
    }, [err])
    let lText = textId
    if (navigator.languages && !navigator.languages.includes('id')) {
        lText = textEn
    }
    if (wasmErr) return <WasmErr />
    if (loading) return <Loading />
    return (
        <div>
            <h1>{lText.title}</h1>
            <h2>{lText.subTitle}</h2>
            <div id='dropper'>
                <div>
                    {err ? (
                        <>
                            <h3>{err}</h3>
                            <h4>{lText.tryAgain}</h4>
                        </>
                    ) : (
                        <>
                            <h3>{lText.dropTitle}</h3>
                            <h4>{lText.dropSubTitle}</h4>
                        </>
                    )}
                </div>
                <input
                    type='file'
                    multiple
                    accept='application/pdf,.pdf'
                    draggable
                    onDragEnter={onDragEnter}
                    onDragOver={onDragOver}
                    onDragEnd={onDragEnd}
                    onDragExit={onDragEnd}
                    onDragExitCapture={onDragEnd}
                    onDragLeave={onDragEnd}
                    onChange={e => {
                        // @ts-ignore
                        const eventFiles = e.target.files as FileList
                        if (eventFiles && eventFiles.length) {
                            let files: File[] = []
                            for (let i = 0; i < eventFiles.length; i++) {
                                const evFile = eventFiles.item(i)
                                if (!evFile) {
                                    continue
                                }
                                if (evFile.type !== 'application/pdf') {
                                    return setErr(lText.errNotPdf)
                                } else {
                                    files.push(evFile)
                                }
                            }
                            processFiles(files)
                            rmActive()
                        } else {
                            rmActive()
                            setErr('')
                        }
                    }}
                    onDrop={e => {
                        const eventFiles = e?.dataTransfer?.files
                        if (eventFiles && eventFiles.length) {
                            let files: File[] = []
                            for (let i = 0; i < eventFiles.length; i++) {
                                const evFile = eventFiles.item(i)
                                if (!evFile) {
                                    continue
                                }
                                if (evFile.type !== 'application/pdf') {
                                    return setErr(lText.errNotPdf)
                                } else {
                                    files.push(evFile)
                                }
                            }
                            processFiles(files)
                            rmActive()
                        } else {
                            rmActive()
                            setErr('')
                        }
                    }}
                />
            </div>
            <p>{lText.remarks}</p>
            <div className='link'>
                <a
                    target='_blank'
                    href='https://github.com/benedictjohannes/bca-pdfestatement-extractor'
                >
                    [Source Code]
                </a>
                <a
                    target='_blank'
                    href='https://www.linkedin.com/in/benedictjohannes/'
                >
                    [LinkedIn]
                </a>
            </div>
        </div>
    )
}

export default App
