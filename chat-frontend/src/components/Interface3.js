import React, {useEffect, useState} from 'react'
import axios from 'axios'
import ComposeModal from "./composeModal";

// Email Detail Component
const EmailDetail = ({email}) => {
    if (!email) {
        return <div>Please select an email</div>
    }

    return (
        <>
            {/*<h2 className="text-lg font-bold mb-2">{email.subject}</h2>*/}
            <div className="text-sm text-gray-500 mb-4">From: {email.sender}</div>
            <p>{email.content}</p>
        </>
    )
}

// Email List Component
const EmailList = ({messages, onEmailSelect}) => {
    return (
        <div className="w-64 bg-white border-r overflow-y-auto">
            {messages.map(message => (
                <div
                    key={message.id}
                    className="p-4 border-b cursor-pointer hover:bg-gray-100"
                    onClick={() => onEmailSelect(message)}
                >
                    <div
                        className={"text-xs text-gray-700 text-right"}>{(new Date(message.timestamp * 1000)).toLocaleString(undefined, {
                        year: 'numeric',    // Display the year
                        month: 'short',     // Display the abbreviated month name
                        day: 'numeric',     // Display the day of the month
                        hour: 'numeric',    // Display the hour (in 12-hour format)
                        minute: '2-digit',  // Display the minute
                        second: '2-digit',  // Display the second
                        hour12: false       // Use 12-hour format
                    })}</div>
                    <div className="font-semibold text-sm">{message.sender.slice(0, 16)}</div>
                    <div className="text-xs text-gray-500">{cutBody(message.content)}</div>
                </div>
            ))}
        </div>
    )
}

// Email Interface Component
const EmailInterface = () => {
    const [selectedEmail, setSelectedEmail] = useState(null)
    const [messages, setMessages] = useState([])
    const [composeModalOpen, setComposeModalOpen] = useState(false)


    useEffect(() => {
        axios.get('http://localhost:8080/messages')
            .then(response => setMessages(response.data))
            .catch(error => console.error('Error fetching emails:', error))
    }, [])

    // Simulate refresh function
    const handleRefresh = () => {
        console.log('Refreshing emails...')
        axios.get('http://localhost:8080/messages')
            .then(response => setMessages(response.data))
            .catch(error => console.error('Error fetching emails:', error))
    }

    const handleCompose = () => {
        setComposeModalOpen(true)
    }

    const handleCloseComposeModal = () => {
        setComposeModalOpen(false)
    }

    const handleSendEmail = ({recipient, body}) => {
        axios.post('http://localhost:8080/send', {"recipient": recipient, "body": body})
            .then(response => {
                console.log('Email sent successfully:', response)
                handleCloseComposeModal()
            })
            .catch(error => console.error('Error sending email:', error))
    }

    return (
        <div className="flex h-screen bg-gray-100">
            <div className="w-64 bg-white border-r overflow-y-auto">
                {/* Buttons */}
                <div className="p-4 flex justify-between border-b">
                    <button className="bg-blue-500 text-white rounded px-4 py-2" onClick={handleRefresh}>Refresh
                    </button>
                    <button className="bg-green-500 text-white rounded px-4 py-2" onClick={handleCompose}>Compose
                    </button>
                </div>

                {/* Email List */}
                <EmailList messages={messages} onEmailSelect={setSelectedEmail}/>
            </div>

            {/* Email Viewer */}
            <div className="flex-1 p-6 overflow-y-auto">
                <EmailDetail email={selectedEmail}/>
            </div>
            <ComposeModal
                isOpen={composeModalOpen}
                onClose={handleCloseComposeModal}
                onSend={handleSendEmail}
            />
        </div>
    )
}

function cutBody(text) {
    // Check if the string exceeds 100 characters
    if (text.length > 100) {
        // Using slice() method to cut off the string at 100 characters
        let slicedStr = text.slice(0, 100);

        // Add an ellipsis at the end of the sliced string
        slicedStr += "...";

        return slicedStr
    } else {
        return text
    }
}

export default EmailInterface
