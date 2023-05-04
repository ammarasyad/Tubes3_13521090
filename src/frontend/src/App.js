import React, { useState } from 'react';
import './App.css';
import ChatRoom from './ChatRoom';

function Sidebar({ chatRooms, activeChatRoom, onSelectChatRoom, onAddChatRoom }) {
  return (
    <div className="sidebar">
      <div className="chat-room-list">
        <ul>
          {chatRooms.map((chatRoom) => (
            <li
              key={chatRoom.id}
              className={chatRoom.id === activeChatRoom ? 'active' : ''}
              onClick={() => onSelectChatRoom(chatRoom.id)}
            >
              {chatRoom.name}
            </li>
          ))}
          <li className="add-chat-room" onClick={onAddChatRoom}>
          Add Chat Room
          </li>
        </ul>
      </div>
    </div>
  );
}


function App() {
  const [activeChatRoom, setActiveChatRoom] = useState(1);
  const [chatHistory, setChatHistory] = useState({});
  const [chatRooms, setChatRooms] = useState([
    { id: 1, name: 'Chat Room 1' },
    { id: 2, name: 'Chat Room 2' },
    { id: 3, name: 'Chat Room 3' },
  ]);

  const handleSelectChatRoom = (chatRoomId) => {
    setActiveChatRoom(chatRoomId);
  };

  const updateChatHistory = (chatRoomId, messages) => {
    setChatHistory((prevChatHistory) => ({
      ...prevChatHistory,
      [chatRoomId]: messages,
    }));
  };

  const getChatHistory = (chatRoomId) => {
    return chatHistory[chatRoomId] || [];
  };

  const handleAddChatRoom = () => {
    const newChatRoom = {
      id: chatRooms.length + 1,
      name: `Chat Room ${chatRooms.length + 1}`,
    };
    setChatRooms([...chatRooms, newChatRoom]);
    setActiveChatRoom(newChatRoom.id);
  };

  return (
    <div className="app">
      <div className="container">
        <Sidebar
          chatRooms={chatRooms}
          activeChatRoom={activeChatRoom}
          onSelectChatRoom={handleSelectChatRoom}
          onAddChatRoom={handleAddChatRoom}
        />
        <div className="main">
          <ChatRoom
            chatRoomId={activeChatRoom}
            chatHistory={getChatHistory(activeChatRoom)}
            updateChatHistory={updateChatHistory}
          />
        </div>
      </div>
    </div>
  );
}

export default App;
