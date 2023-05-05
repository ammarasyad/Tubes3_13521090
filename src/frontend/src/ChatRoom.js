import React, { useState, useEffect } from 'react';
import './ChatRoom.css';

function ChatRoom({ chatRoomId, chatHistory, updateChatHistory }) {
  const [messages, setMessages] = useState([]);
  const [currentMessage, setCurrentMessage] = useState('');
  const [selectedOption, setSelectedOption] = useState('option1');
  const handleOptionChange = (event) => {
    setSelectedOption(event.target.value);
  };

  useEffect(() => {
    // Load chat history for the current chat room
    setMessages(chatHistory);
  }, [chatHistory]);

  useEffect(() => {
    const initHistory = async () => {
      if (chatRoomId === 1) {
        fetch('http://localhost:3001/api/get/history')
          .then((response) => response.json())
          .then((data) => {
            console.log(data);
            const newResponseArray = data.questions.reverse().map((string, index) => ([{
              sender: 1,
              id: messages.length + 1 + index * 2,
              content: string,
            }, {
              sender: 0,
              id: messages.length + 2 + index * 2,
              content: data.answers[data.answers.length - 1 - index],
            }]));
            const updatedMessages = [...messages, ...newResponseArray.flat()];
            console.log(newResponseArray.flat());
            setMessages(updatedMessages);
            updateChatHistory(chatRoomId, updatedMessages); // Update the chat history in App
            // Handle the response data if needed
          })
          .catch((error) => {
            console.error('Error:', error);
          });
      }
    };
  
    initHistory();
  }, [chatRoomId]);

  const handleSendMessage = () => {
    if (currentMessage.trim() !== '') {
      const newMessage = {
        sender: 1,
        id: messages.length + 1,
        content: currentMessage,
      };

      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: currentMessage,
                                kmpbm: selectedOption === 'option1'}),
      };
  
      fetch('http://localhost:3001/api/messages', requestOptions)
        .then(response => response.json())
        .then(data => {
          console.log(data);
          const newResponse = {
            sender: 0,
            id: messages.length + 2,
            content: data.result,
          };
          const updatedMessages = [...messages, newMessage, newResponse];
          setMessages(updatedMessages);
          updateChatHistory(chatRoomId, updatedMessages); // Update the chat history in App
          // Handle the response data if needed
          const saveOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ question: currentMessage,
                                    answer: data.result}),
          };
          return fetch('http://localhost:3001/api/save/history', saveOptions);
        })
        .then(response => response.text())
        .then(data => {
          console.log(data);
        })
        .catch(error => {
          console.error('Error:', error);
        });

      setCurrentMessage('');
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const formatMessageContent = (content) => {
    return content.split('\n').map((line, index) => (
      <React.Fragment key={index}>
        {line}
        <br />
      </React.Fragment>
    ));
  };

  return (
    <div className="chat-room">
      <div className="chat-box">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`message-bubble-${message.sender === 1 ? 'right' : 'left'}`}
          >
            {formatMessageContent(message.content)}
          </div>
        ))}
      </div>
      <div className="message-box">
        <textarea
          className="message-input"
          value={currentMessage}
          onChange={(e) => setCurrentMessage(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type a message"
          rows={1}
          style={{ resize: 'none' }}
        />
        <button className="send-button" onClick={handleSendMessage}>
          Send
        </button>
        <label>
        <input
          type="radio"
          value="option1"
          checked={selectedOption === 'option1'}
          onChange={handleOptionChange}
        />
        KMP
      </label>

      <label>
        <input
          type="radio"
          value="option2"
          checked={selectedOption === 'option2'}
          onChange={handleOptionChange}
        />
        BM
      </label>
      </div>
    </div>
  );
}

export default ChatRoom;
