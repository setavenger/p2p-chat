import React, { useState } from 'react'

// ComposeModal component
const ComposeModal = ({ isOpen, onClose, onSend, message }) => {
    const [recipient, setRecipient] = useState('')
    const [body, setBody] = useState('')

    const handleSubmit = e => {
        e.preventDefault()
        onSend({ recipient, body })
        setRecipient("")
        setBody("")
    }

    if (!isOpen) {
        return null
    }

    return (
        <div className="fixed top-0 left-0 w-full h-full flex items-center justify-center bg-black bg-opacity-50">
            <form className="bg-white p-8 rounded shadow-md w-1/4" onSubmit={handleSubmit}>
                <h2 className="text-xl font-bold mb-4 w-96">Compose Reply Message</h2>
                {/*<div className="block mb-2">*/}
                {/*    Reply to:*/}
                {/*</div>*/}
                <label className="block mb-4">
                    Message:
                    <textarea
                        className="border rounded w-full py-2 px-4"
                        value={body}
                        placeholder={"Reply comment..."}
                        onChange={e => setBody(e.target.value)}
                    />
                </label>
                <div className="flex justify-end">
                    <button type="button" className="mr-2" onClick={onClose}>Cancel</button>
                    <button type="submit" className="bg-blue-500 text-white rounded px-4 py-2">Send</button>
                </div>
            </form>
        </div>
    )
}

export default ComposeModal