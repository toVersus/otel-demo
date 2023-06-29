// import require from './require'
// require("dotenv").config()

// Components
const Title = (props) => <p>{props.label}</p>
const Response = (props) => <pre>{JSON.stringify(props.result, null, 4)}</pre>
const TransferAmount = (props) => {
  return (
    <select onChange={props.handleTransferAmountChange}>
      <option selected>Choose amount</option>
      <option value={500}>500</option>
      <option value={800}>800</option>
      <option value={1500}>1500</option>
    </select>
  )
}
const ProductDD = (props) => {
  return (
    <select onChange={props.handleProductChange}>
      <option selected>Choose product</option>
      <option value={599}>Mac Mini</option>
      <option value={999}>MacBook Air</option>
      <option value={1100}>MacBook Pro</option>
    </select>
  )
}

const App = () => {
  let [user, setUser] = React.useState()
  let [payment, setPayment] = React.useState()
  let [order, setOrder] = React.useState()
  let [amount, setAmount] = React.useState()
  let [product, setProduct] = React.useState()

  let createUser = async () => {
    try {
      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_name: 'Naman',
          account: 'savings',
        }),
      }
      let userServiceUrl = '/users';
      if (typeof USER_URL !== 'undefined') {
        userServiceUrl = USER_URL + userServiceUrl
      }
      let response = await fetch(userServiceUrl, requestOptions)
      let result = await response.json()
      setUser(result)
      console.log(result)
    } catch (error) {
      setUser(error.message)
      console.log('err', error.message)
    }
  }

  let handleProduct = (e) => {
    let { options, value } = e.target
    setProduct({
      product_name: options[options.selectedIndex].text,
      price: value,
    })
  }

  let handleAmount = (e) => {
    let { _, value } = e.target
    setAmount({
      value: Number(value),
    })
  }

  let transferFund = async () => {
    console.log('transferFund', user)
    console.log('amount', amount)
    console.log(typeof amount.value)
    const requestOptions = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        amount: amount.value,
      }),
    }

    let paymentServiceUrl = `/payments/transfer/id/${user.id}`;
    if (typeof PAYMENT_URL !== 'undefined') {
      paymentServiceUrl = PAYMENT_URL + paymentServiceUrl
    }
    console.log('paymentServiceUrl', paymentServiceUrl)
    let response = await fetch(
      paymentServiceUrl,
      requestOptions
    )
    let result = await response.json()
    setPayment(result)
    console.log(result)
  }

  let placeOrder = async () => {
    console.log('placeOrder', product)
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: user.id,
        product_name: product.product_name,
        price: parseInt(product.price),
      }),
    }

    let orderServiceUrl = '/orders';
    if (typeof ORDER_URL !== 'undefined') {
      orderServiceUrl = ORDER_URL + orderServiceUrl
    }
    let response = await fetch(orderServiceUrl, requestOptions)
    let result = await response.json()
    setOrder(result)
    console.log(result)
  }

  let handleReset = () => {
    setUser()
    setPayment()
    setAmount()
    setOrder()
    setProduct()
  }

  return (
    <div>
      <button onClick={handleReset}>Reset Actions</button>
      <Title label="1. User Creation"></Title>
      <button onClick={createUser}>Create User</button>
      <Response result={user} />
      {user && (
        <div>
          <Title label="2. Transfer amount"></Title>
          <TransferAmount handleTransferAmountChange={handleAmount} />
          {amount && <button onClick={transferFund}>Transfer Fund</button>}
          <Response result={payment} />
        </div>
      )}
      {payment && (
        <div>
          <Title label="3. Place order"></Title>
          <ProductDD handleProductChange={handleProduct} />
          {product && <button onClick={placeOrder}>Place Order</button>}
          {order && (
            <div>
              <Response result={order} />
              <h3>Order Placed!</h3>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

ReactDOM.render(<App />, document.getElementById('app'))
