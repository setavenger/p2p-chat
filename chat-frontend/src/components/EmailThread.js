import React from 'react';


const Message = ({ message, level = 0 }) => {
    return (
        <div className={`mt-2 mb-2 ${level === 0 ? '' : `ml-${level*4} pl-2 border-l-2 border-blue-500`}`}>
            <div>{message.text}</div>
            {message.replies.map(reply =>
                <Message key={reply.id} message={reply} level={level + 1} />
            )}
        </div>
    );
}

export const EmailThread = ({messages}) => {
    return (
        <div className={"whitespace-pre-line"}>
            {messages.map(message =>
                <Message key={message.id} message={message}/>
            )}
        </div>
    );
}
