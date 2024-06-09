import logo from './logo.svg';
import React, { useState } from 'react';
import './App.css';

function App() {
  const [orderAmount, setOrderAmount] = useState('');
  const [packageSizes, setPackageSizes] = useState(Array(5).fill(''));
  const [responseMessage, setResponseMessage] = useState('');
  const [orderHistory, setOrderHistory] = useState([]);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginError, setLoginError] = useState('');


  const handlePackageSizeChange = (index, value) => {
    setPackageSizes(packageSizes.map((size, i) => (i === index ? value : size)));
  };

  const handleLogin = async () => {
    try {
      const response = await fetch('http://localhost:8080/loginUser', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'valid-token'
        },
        body: JSON.stringify({
          username: username,
          password: password,
        }),
      });

      if (response.ok) {
        setIsAuthenticated(true);
        setLoginError('');
      } else {
        throw new Error('Failed to authenticate');
      }
    } catch (error) {
      setLoginError(error.message);
    }
  };

  const sendPostRequest = async () => {
    try {
      const response = await fetch('http://localhost:8080/order', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'valid-token'
        },
        body: JSON.stringify({
          orderAmount: Number(orderAmount),
          packageSizes: packageSizes.map(Number)
        })
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      if (response.status !== 204) {
        const data = await response.json();
        setResponseMessage(`Amount: ${data.amount}, Result: ${data.result.join(', ')}`); // Store the response in the state
      } else {
        setResponseMessage('POST request was successful but no content returned');
      }

      const historyResponse = await fetch('http://localhost:8080/getDocument', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'valid-token'
        }
      });

      if (!historyResponse.ok) {
        throw new Error(`HTTP error! status: ${historyResponse.status}`);
      }

      const historyData = await historyResponse.json();
      setOrderHistory(historyData.history); // Store the history in the state
    } catch (error) {
      console.error('Failed to fetch:', error);
    }
  };

  return (
    <div className="App">
      {!isAuthenticated ? (
        <div className="login-modal">
          <div className="login-popup">
            {loginError && <p className="login-error">{loginError}</p>}
            <input
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              placeholder="Username"
            />
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              placeholder="Password"
            />
            <button onClick={handleLogin}>Login</button>
          </div>
        </div>
      ) : (
        <>
          <header className="App-header">
            <div className="package-sizes">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="package-size">
                  <label htmlFor={`package-size-${i + 1}`} className="package-size-label">Package {i + 1}</label>
                  <input
                    id={`package-size-${i + 1}`}
                    type="number"
                    value={packageSizes[i]}
                    onChange={e => handlePackageSizeChange(i, e.target.value)}
                    placeholder={`Size ${i + 1}`}
                    className="package-size-input"
                  />
                </div>
              ))}
            </div>
            <div className="order-amount">
              <label htmlFor="order-amount" className="order-amount-label">Order Amount</label>
              <input
                id="order-amount"
                type="number"
                value={orderAmount}
                onChange={e => setOrderAmount(e.target.value)}
                placeholder="Amount"
                className="order-amount-input"
              />
            </div>
            <button onClick={sendPostRequest} className="submit-button">Submit</button>
          </header>
          <div className="response-box">{responseMessage}</div>
          <h2 className="order-history-title">Order History</h2>
          <div className="order-history">
            {orderHistory.map((order, i) => (
              <div key={i} className="order-history-item">
                <p>Order {i + 1}: Amount: {order.order.amount}, Result: {order.order.result.join(', ')}</p>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}

export default App;