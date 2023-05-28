import React, {useEffect, useState} from 'react'
import axios from 'axios'
import ComposeModal from "./composeModal";
import SetKeyModal from "./setKeyModal";
import SetHostModal from "./setHostModal";
import {baseUrl, changePublicKey, publicKey} from "../constants";
import ReplyModal from "./replyModal";


const MessageDetail = ({selectedMessage, messages}) => {
    const renderMessage = (message, level) => {
        return (
            <div>
                <div>{message.content}</div>
                {message.parent_id && renderParentMessage(message.parent_id)}
            </div>
        );
    };

    const renderParentMessage = (parentId, level = 1) => {
        const parentMessage = messages.find((message) => message.id === parentId);
        if (parentMessage) {
            return (
                <div
                    className={`pl-${level * 4} border-l-2 ${parentMessage.sender === publicKey ? "border-gray-500" : "border-amber-500"}`}>
                    <div className={"mt-2"}>
                        {renderMessage(parentMessage, level + 1)}
                    </div>
                </div>
            );
        }
        return null;
    };

    return (
        <div>
            {selectedMessage && renderMessage(selectedMessage)}
        </div>
    );
};


// Email List Component
const EmailList = ({messages, onEmailSelect}) => {
    return (
        <div className="w-64 bg-white overflow-y-auto">
            {messages !== null ? messages.map(message => (
                message.sender !== publicKey ?
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
                        <div className="font-semibold text-sm">{message.sender_username}</div>
                        <div className="text-xs text-gray-500">{cutBody(message.content)}</div>
                    </div> : null
            )) : <div className={"text-center mt-10"}>No messages yet</div>}
        </div>
    )
}

// Email Interface Component
const EmailInterface = () => {
    const [selectedMessage, setSelectedMessage] = useState(null)
    const [messages, setMessages] = useState([])
    const [composeModalOpen, setComposeModalOpen] = useState(false)
    const [keyModalOpen, setKeyModalOpen] = useState(false)
    const [hostModalOpen, setHostModalOpen] = useState(false)
    const [replyModalOpen, setReplyModalOpen] = useState(false)


    useEffect(() => {
        axios.get(baseUrl + '/messages')
            .then(response => {
                console.log(response.data)
                setMessages(response.data)
            })
            .catch(error => console.error('Error fetching emails:', error))
    }, [])

    // get publicKey info
    useEffect(() => {
        axios.get(baseUrl + '/meta')
            .then(response => {
                changePublicKey(response.data.public_key)
            })
            .catch(error => console.error('Error fetching meta data:', error))
    }, [])

    // Simulate refresh function
    const handleRefresh = () => {
        console.log('Refreshing emails...')
        axios.get(baseUrl + '/messages')
            .then(response => setMessages(response.data))
            .catch(error => console.error('Error fetching emails:', error))
    }

    const handleComposeModal = () => {
        setComposeModalOpen(true)
    }

    const handleCloseComposeModal = () => {
        setComposeModalOpen(false)
    }

    const handleKeyModal = () => {
        setKeyModalOpen(true)
    }

    const handleCloseKeyModal = () => {
        setKeyModalOpen(false)
    }

    const handleHostModal = () => {
        setHostModalOpen(true)
    }

    const handleCloseHostModal = () => {
        setHostModalOpen(false)
    }
    const handleReplyModal = () => {
        setReplyModalOpen(true)
    }

    const handleCloseReplyModal = () => {
        setReplyModalOpen(false)
    }

    const handleSendMessage = ({recipient, body}) => {
        axios.post(baseUrl + '/send', {"recipient": recipient, "body": body})
            .then(response => {
                console.log('Email sent successfully:', response)
                handleCloseComposeModal()
            })
            .catch(error => console.error('Error sending email:', error))
    }
    const handleSendReply = ({recipient, body}) => {
        if (selectedMessage === null) {
            console.log("No message was selected")
            return
        }
        axios.post(baseUrl + '/send', {"recipient": selectedMessage.sender_username, "body": body, "parent_id": selectedMessage.id})
            .then(response => {
                console.log('Email sent successfully:', response)
                handleCloseComposeModal()
            })
            .catch(error => console.error('Error sending email:', error))
    }

    const handleSetKey = ({key}) => {
        axios.post(baseUrl + '/set-key', {"key": key})
            .then(response => {
                console.log('key set successfully:', response)
                handleCloseKeyModal()
                handleRefresh()
            })
            .catch(error => console.error('Error setting key:', error))
    }

    const handleSetHost = ({host}) => {
        axios.post(baseUrl + '/set-host', {"host": host})
            .then(response => {
                console.log('host set successfully:', response)
                handleCloseHostModal()
                handleRefresh()
            })
            .catch(error => console.error('Error setting host:', error))
    }

    return (
        <>
            <div className={"flex h-screen bg-gray-100"}>
                <div className={"border-r h-full bg-white"}>
                    <EmailList messages={messages} onEmailSelect={setSelectedMessage}/>
                </div>
                <div className={"flex-1 flex-col"}>
                    <div className="p-4 flex justify-between ">
                        <button className="bg-blue-500 text-white rounded px-4 py-2" onClick={handleRefresh}>
                            Refresh
                        </button>
                        <button className="bg-gray-500 text-white rounded px-4 py-2" onClick={handleReplyModal}>
                            Reply
                        </button>
                        <button className="bg-green-500 text-white rounded px-4 py-2" onClick={handleComposeModal}>
                            Compose
                        </button>
                        <button className="bg-red-500 text-white rounded px-4 py-2" onClick={handleKeyModal}>
                            Set Key
                        </button>
                        <button className="bg-yellow-500 text-white rounded px-4 py-2" onClick={handleHostModal}>
                            Set Host
                        </button>
                    </div>
                    <div className="flex-1 p-6 overflow-y-auto">
                        <MessageDetail selectedMessage={selectedMessage} messages={messages}/>
                    </div>
                </div>
                <ComposeModal
                    isOpen={composeModalOpen}
                    onClose={handleCloseComposeModal}
                    onSend={handleSendMessage}
                />
                <ReplyModal
                    selectedMessage={selectedMessage}
                    isOpen={replyModalOpen}
                    onClose={handleCloseReplyModal}
                    onSend={handleSendReply}
                />
                <SetKeyModal
                    isOpen={keyModalOpen}
                    onClose={handleCloseKeyModal}
                    onSend={handleSetKey}
                />
                <SetHostModal
                    isOpen={hostModalOpen}
                    onClose={handleCloseHostModal}
                    onSend={handleSetHost}
                />
            </div>
        </>
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
