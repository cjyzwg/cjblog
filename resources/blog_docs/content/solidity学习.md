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

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;
contract MappingExample{
    mapping (address => uint) account;
    function get(address _address) public view returns (uint){
        return account[_address];
    }
    function set(address _address,uint _balance) public {
        account[_address] = _balance;
    }
    function remove(address _address) public {
        delete  account[_address];
    }
}

```
提到了mapping嵌套，具体还得再查下，但是mapping嵌套我理解上是二维数组类型的。

### ERC20代币的介绍

使用openZeppelin 实现

### 数组应用

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;
contract ArrayExample{
    uint[] iArray;
    uint[] iArray2 = [1,2,3];
    uint[3] iArray3;
    function getArray() public  view returns (uint[] memory){
        return iArray2;
    }
    function getArrayByIndex(uint _i) public  view returns (uint) {
        return iArray2[_i];
    }
    function getLength() public view returns (uint){
        return iArray3.length;
    }
    function push(uint _i) public {
        iArray2.push(_i);
    }
    function pop() public {
        iArray2.pop();
    }
    function deleteByIndex(uint _i) public {
        delete iArray2[_i];
    }
}

```


### imort引入文件应用

```ts

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "./13_import_1.sol";
//引入也可以引入外部链接比如
//import "https://github.com/13_import_1.sol";
contract ImportExample2{
    ImportExample importExample = new ImportExample();
    function getAge() public view returns (uint) {
        return importExample.age();
    }
    function getName() public view returns (string memory){
        return importExample.getName();
    }
}

```


### 创建最简单DEX，进行ERC20交易

```ts

// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.2;
library SafeMath{
    function sub(uint256 a,uint256 b) internal pure returns (uint256){
        assert(b%3C=a);
        return a-b;
    }
    function add(uint256 a,uint256 b) internal pure returns (uint256){
        uint256 c = a+b;
        assert(c%3E=a);
        return c;
    }
}
interface IERC20 {
    function getAddress() external view returns(address);
    function totalSupply() external view returns(uint256);
    function balanceOf(address account) external view returns (uint256);
    function allowance(address owner,address spender) external view returns (uint256);
    
    function transfer(address recipient,uint256 amount) external returns (bool);
    function approve(address owner,address spender,uint256 amount) external returns (bool);
    function transferFrom(address sender,address recipient,uint256 amount) external returns(bool);

    event Transfer(address indexed from,address indexed to,uint256 value);
    event Approval(address indexed owner,address indexed spender,uint256 value);
}

contract ERC20Basic is IERC20 {

    string public constant name = "ERC20-ThinkChain";
    string public  constant symbol = "ERC-TKC";
    uint8 public constant decimals = 18;
    mapping(address=>uint256) balances;
    mapping(address=>mapping(address=>uint256)) allowed;

    uint256 totalSupply_ = 10 ether;
    using SafeMath for uint256;
    constructor() {
        balances[msg.sender] = totalSupply_;
    }
    function getAddress() public  override view returns (address){
        return address(this);
    }
    function totalSupply() public override  view returns (uint256){
        return totalSupply_;
    }
    function balanceOf(address tokenOwner) public  override view returns (uint256) {
        return balances[tokenOwner];
    }
    function transfer(address receiver,uint256 numTokens) public override returns (bool){
        require(numTokens <= balances[msg.sender]);
        balances[msg.sender] = balances[msg.sender].sub(numTokens);
        balances[receiver] = balances[receiver].add(numTokens);
        emit Transfer(msg.sender,receiver,numTokens);
        return true;
    }
    function approve(address owner,address delegate,uint256 numTokens) public  override returns (bool){
        allowed[owner][delegate] = numTokens;
        emit Approval(owner, delegate, numTokens);
        return true;
    }
    function allowance(address owner,address delegate) public override view returns (uint){
        return allowed[owner][delegate];
    }
    function transferFrom(address owner,address buyer,uint256 numTokens) public override returns (bool){
        require(numTokens<=balances[owner]);
        require(numTokens<=allowed[owner][msg.sender]);
        balances[owner] = balances[owner].sub(numTokens);
        allowed[owner][buyer] = allowed[owner][buyer].sub(numTokens);
        emit Transfer(owner, buyer, numTokens);
        return true;
    }
}

contract DEX{
    event Bought(uint256 amount);
    event Sold(uint256 amount);
    IERC20 public token;
    constructor(){
        token = new ERC20Basic();
    }
    function buy() payable public {
        uint256 amountTobuy = msg.value;
        uint256 dexBalance = token.balanceOf(address(this));
        require(amountTobuy>0,"You need to send some Ether");
        require(amountTobuy<=dexBalance,"Not enough tokens in the reserve");
        token.transfer(msg.sender, amountTobuy);
        emit Bought(amountTobuy);
    }
    function sell(uint256 amount) public {
        require(amount>0,"You need to sell at least some tokens");
        uint256 allowance = token.allowance(msg.sender, address(this));
        require(allowance>=amount,"Check the token allowance");
        token.transferFrom(msg.sender, address(this), amount);
        emit  Sold(amount);
    }
    function getDexBalance() public view returns (uint256) {
        return token.balanceOf(address(this));
    }
    function getOwnerBalance() public view returns (uint256) {
        return token.balanceOf(msg.sender);
    }
}

```

### openzeppelin 合约

铸币合约

```ts

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
contract ThinkingChainToken is ERC20 {
    constructor() ERC20("ThinkingChain","TKC") {
        _mint(msg.sender, 1000 * 10 ** decimals());
    }
}

```

接受代币的合约

```ts

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
contract ReciveTokenContract{

    IERC20 myToken;
    constructor(address _tokenAddress) {
        myToken = IERC20(_tokenAddress);
    }
    function transferFrom(uint _amount) public {
        myToken.transferFrom(msg.sender,address(this),_amount);
    }

    function getBalance(address _address) public  view returns(uint) {
        return myToken.balanceOf(_address);
    }
}

```

### interface 接口应用

先部署Employee 合约 再部署company合约
```ts

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
interface IEmployee {
    
    function setName(string memory _name) external ;
    function getName() external view returns (string memory);

}

contract Employee is IEmployee{
    string private  name;
    function setName(string memory _name) public override {
        name = _name;
    }
    function getName() public override view returns (string memory){
        return name;
    }

}

contract Company{
    IEmployee employee;
    constructor(address _address) {
        employee = IEmployee(_address);
    }
    function setName(string memory _name) public  {
        employee.setName(_name);
    }
    function getName() public  view returns (string memory){
        return employee.getName();
    }
}

```

### Library应用

使用类库的时候是_a.add(_b) 不是用safemath.add(_a,_b)

```ts

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
library SafeMath {
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");

        return c;
    }

    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b > 0, "SafeMath: division by zero");
        uint256 c = a / b;
        return c;
    }

    function sub(uint256 a, uint256 b) public  pure returns (uint256) {
        require(b <= a, "SafeMath: subtraction overflow");
        uint256 c = a - b;

        return c;
    }
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "SafeMath: addition overflow");

        return c;
    }
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0, "SafeMath: modulo by zero");
        return a % b;
    }
}

contract Example{
    using SafeMath for uint256;
    function doAdd(uint256 _a,uint256 _b) public pure returns(uint256){
        return _a.add(_b);
    }
    function doSub(uint256 _a,uint256 _b) public pure returns(uint256){
        return _a.sub(_b);
    }
    function doAddMore(uint256 _a,uint256 _b,uint256 _c) public pure returns (uint256) {

return _a.addMore(_b,_c);

}
}

```
### 代理合约

 用于智能合约的升级

```ts
//logic.sol

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
contract Logic {
    uint private number;
    function setNumber(uint _number) public {
        number = _number+1;
    }
    function getNumber() public view returns (uint){
        return number;
    }
}

```


```ts
//logic2.sol

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
contract Logic2 {
    uint private number;
    function setNumber(uint _number) public {
        number = _number+2;
    }
    function getNumber() public view returns (uint){
        return number;
    }
}

```

```ts
//Proxy.sol 代理合约设置借口

// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
interface ILogic {
    function setNumber(uint _number) external ;
    function getNumber() external view returns(uint); 
}

contract Proxy {
    ILogic public logic;
    function setLogicAddress(address _logicAddress) public {
        logic = ILogic(_logicAddress);
    }
    function getLogicAddress() public view returns (address){
        return address(logic);
    }
    function setNumber(uint _number) public {
        logic.setNumber(_number);
    }
    function getNumber() public view returns (uint){
        return logic.getNumber();
    }
}

```

