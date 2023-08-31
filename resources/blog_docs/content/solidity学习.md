---
title: solidity学习
tags: [solidity,博客,'#Fleeting_N']
date: 2023-08-31 19:13:18
draft: true
hideInList: false
isTop: false
published: true
categories: [solidity,博客]
---

作者：阿三 
博客：[Nockygo](https://hexo.hexiefamily.xin) 
公众号：阿三爱吃瓜

> 持续不断记录、整理、分享，让自己和他人一起成长！😊


------

### modifier作用

其实就是可以作为校验器来处理，比如限制权限，只能在合约部署的人才可以使用，其他账号只能做其他处理，或者金额要超过多少。

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ExampleModifier{
    address public owner;
    uint256 public account;
    constructor(){
        owner = msg.sender;
        account = 0;
    }
    modifier onlyOwner(uint256 _account){
        require(msg.sender == owner,"Only Owner");
        require(_account>100,"Valid 100");
        _;
    }

    function updateAccount(uint256 _account) public onlyOwner(_account){
        account = _account;

    }
}

```

### event 作用

其实就是event emit 提交 写log日志

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ExampleEvent{
    event Deposit(address _from,string _name,uint256 _value);
    function deposit(string memory _name) public payable {
        emit Deposit(msg.sender, _name, msg.value);
    }
}

```

### view pure区别

pure 不允许访问状态变量 也不允许更改
view 允许访问 不允许更改
```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ExampleView{
    uint public age;
    function increaseAge() public {
        age ++;
    }
    function getAgeWithView() public view returns(uint){
        
        return age;
    }
    function getAgeWithPure(uint _age) public pure returns (uint){
        _age ++;
        return _age;
    }
}

```

### public private internal external 用法

external 只能从外部访问

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract Salary{
    uint public data;
    function getData() external view returns (uint) {
        return data;
    }
    function setData(uint _data) external {
        data = _data;
    }
}
//只会部署这个合约
contract Employee {
    Salary salary;
    constructor(){
        salary = new Salary();
    }
    function getSalary() external view returns (uint) {
        return salary.getData();
    }
    function setSalary(uint _data) external {
        salary.setData(_data);
    }
}

```

### address 有哪些
- 合约地址，创建不会改变
- owner地址
- 与合约打交道的人的地址

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ExampleAddress{
    address ownerAddress;
    constructor() {
        ownerAddress = msg.sender;
    }
    function getContractAddress() external view returns (address){
        return address(this);
    }

    function getSenderAddress() external view returns (address){
        return address(msg.sender);
    }

    function getOwnerAddress() external  view  returns (address){
        return address(ownerAddress);
    }
	function getBalance() external view returns (uint,uint,uint){
        uint contractBalance = address(this).balance;
        uint senderBalance = msg.sender.balance;
        uint ownerBalance = ownerAddress.balance;
        return (contractBalance,senderBalance,ownerBalance);
    }
}

```

### address支付方法

接收和转账代币参数和return 必须追加payable

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ExampleAddressFunction{
    function send(address payable _to) public  payable {
        bool isSend = _to.send(msg.value);
        require(isSend,"Send fail");
    }
    function transfer(address payable _to) public  payable {
        _to.transfer(msg.value);
    }
//官方推荐使用call方法
    function call(address payable _to) public  payable {
        (bool isSend,) = _to.call{value:msg.value,gas:5000}("");
        require(isSend,"Send fail");
    } 
}

```

### call 方法 支付+调用智能合约

正常call 方法后面不加参数，会进入到receive方法，
如果存在参数，就会进入fallback 方法，同时返回的log日志中data是空的，
如果想要data中数据不为空，那就重新写一个方法，使用
```ts

function foo(string memory _message,uint _y) public payable returns (uint){
        emit Received(msg.sender, msg.value, _message);
        return _y;
}
abi.encodeWithSignature("foo(string,uint256)", "call foo",_y)

```

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract Receiver {
    event Received(address caller,uint amount,string message);
    receive() external payable {
        emit Received(msg.sender,msg.value,"Receive was called");
    }
    fallback() external  payable {
        emit Received(msg.sender,msg.value,"Fallback was called");
    }
    function foo(string memory _message,uint _y) public payable returns (uint){
        emit Received(msg.sender, msg.value, _message);
        return _y;
    }

    function getAddress() public view returns(address){
        return address(this);
    }
    function getBalance() public view returns (uint) {
        return address(this).balance;
    }
}

contract Caller {
    Receiver receiver;
    constructor(){
        receiver = new Receiver();
    }
    event Response(bool success,bytes data);

    function testCall(address payable _addr,uint _y) public payable {
        (bool success,bytes memory data) = _addr.call{value:msg.value}(
            abi.encodeWithSignature("foo(string,uint256)", "call foo",_y)
        );
        emit Response(success, data);

    }
    function getAddress() public  view returns (address){
        return receiver.getAddress();
    }
    function getBalance() public  view returns (uint) {
        return receiver.getBalance();
    }
}

```

### constant 和Immutable区别

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ConstantImmutable{
    string constant name = "Biden";
    uint immutable age;
    constructor() {
        age = 100;
    }
    //获取constant变量方法必须使用pure修饰符，view 修饰符是错误的
    function getName() public pure returns (string memory){
        return name;
    }
    //如果constructor中未定义age，直接在immutable age=100,那么也必须使用pure修饰符
    function getAge() public view  returns (uint) {
        return age;
    }
}

```

### mapping 的应用场景

```ts

```
