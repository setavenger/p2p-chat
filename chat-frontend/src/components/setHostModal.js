import React, { useState } from 'react'

// ComposeModal component
const SetKeyModal = ({ isOpen, onClose, onSend }) => {
    const [host, setHost] = useState('')

    const handleSubmit = e => {
        e.preventDefault()
        onSend({ host })
        setHost("")
    }

    if (!isOpen) {
        return null
    }

    return (
        <div className="fixed top-0 left-0 w-full h-full flex items-center justify-center bg-black bg-opacity-50">
            <form className="bg-white p-8 rounded shadow-md w-1/4" onSubmit={handleSubmit}>
                <h2 className="text-xl font-bold mb-4 w-96">Set Host</h2>
                <label className="block mb-2">
                    Host address:
                    <input
                        className="border rounded w-full py-2 px-4"
                        type="text"
                        value={host}
                        placeholder={"https://p2p.domain.com"}
                        onChange={e => setHost(e.target.value)}
                    />
                </label>
                <div className="flex justify-end">
                    <button type="button" className="mr-2" onClick={onClose}>Cancel</button>
                    <button type="submit" className="bg-blue-500 text-white rounded px-4 py-2">Set</button>
                </div>
            </form>
        </div>
    )
}

export default SetKeyModal