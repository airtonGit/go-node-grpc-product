const grpc = require('@grpc/grpc-js')
const protoLoader = require('@grpc/proto-loader')

const LISTEN_ADDR = process.env.DISCOUNT_LISTEN_ADDR || "127.0.0.1:50051";
const BIRTHDAY_DISCOUNT = process.env.DISCOUNT_BIRTHDAY_DISCOUNT || 5.00; //5%
const BALCKFRIDAY_DISCOUNT = process.env.DISCOUNT_BLACKFRIDAY_DISCOUNT || 10.00; //10%
const BALCKFRIDAY_DATE = process.env.DISCOUNT_BLACKFRIDAY_DATE || "11/25"

var PROTO_PATH = __dirname + '/../product.proto';

const packageDefinition = protoLoader.loadSync(PROTO_PATH);
const product_proto = grpc.loadPackageDefinition(packageDefinition);

const checkSameMonthDay = function(date1Str, date2Date){
    
    let dateNoYearRegex = /(\d{1,2})\/+?(\d{1,2})/

    let regexObj = new RegExp(dateNoYearRegex);
    
    if (!regexObj.test(date1Str) ){
        console.error(`Date information wrong format given date1Str: ${date1Str} want mm/dd two digits month/day`)
        throw Error(`Date information wrong format given date1Str: ${date1Str} want mm/dd two digits month/day`)
    }
   
    let found = date1Str.match(dateNoYearRegex)

    return found[1] == date2Date.getMonth()+1 && found[2] == date2Date.getDate()
}

const user_list = [
    {id: '1', date_of_birth: '10/21/1983'},
    {id: '2', date_of_birth: '09/03/1987'},
    {id: '3', date_of_birth: '10/20/1987'}
]

const product_list = [
    { id: '1', price_in_cents: 330 * 100 },
    { id: '2', price_in_cents: 299 * 100 },
    { id: '3', price_in_cents: 279 * 100 },
]


const userMonthDayOfBirth = (userFullBirthDayStr) => {
    let dateMonthDay = /(\d{1,2}\/+?\d{1,2})/

    let regexObj = new RegExp(dateMonthDay);
    
    if (!regexObj.test(userFullBirthDayStr) ){
        console.error(`Date information wrong format given userFullBirthDayStr: ${userFullBirthDayStr} want mm/dd two digits month/day`)
        throw Error(`Date information wrong format given userFullBirthDayStr: ${userFullBirthDayStr} want mm/dd two digits month/day`)
    }
    return userFullBirthDayStr.match(dateMonthDay)[1]
}

const getUser = (userID) => {
    let user =  user_list.find( element => element.id === userID )
    if (user == undefined){
        throw Error("getUser user not found");
    }
    return user;
}

const getProduct = (productID) => {
    console.log("getProduct", productID)
    let product =  product_list.find( element => element.id === productID )
    if (product == undefined){
        throw Error("getproduct product not found");
    }
    return product;
}

const checkUserBirthday = (userID, todayDate) => {
    return checkSameMonthDay(userMonthDayOfBirth(getUser(userID).date_of_birth), todayDate);
}

const findDiscount = (user, product) => {
    //check today BlackFriday
        //birtday give discount
    //no birtday try product discount
    let returnProduct = getProduct(product)

    //check Black Friday
    let isBlackfriday = checkSameMonthDay(BALCKFRIDAY_DATE, new Date())

    if (isBlackfriday){
        returnProduct.discount = {
            pct: BALCKFRIDAY_DISCOUNT,
            value_in_cents:  (BALCKFRIDAY_DISCOUNT / 100) * returnProduct.price_in_cents
        }
    }

    //no Blackfriday, try birthday
    if (!isBlackfriday && checkUserBirthday(user, new Date()) ){
        returnProduct.discount = {
            pct: BIRTHDAY_DISCOUNT,
            value_in_cents:  (BIRTHDAY_DISCOUNT / 100) * returnProduct.price_in_cents
        }
    }

    return returnProduct
}

//Implements Discount RPC method
function discount(call, callback){
    //console.log("discount request", call.request)

    let product = {}

    try{
        product = findDiscount(call.request.userId, call.request.productId)
    }catch(err){
        console.error("fail calculate discount err:", err)
    }

    console.log("product with discount", product)

    var response = {
        pct: product.discount.pct,
        valueInCents: Math.floor(product.discount.value_in_cents)
    }

    console.log("response", response)
    
    callback(null, response);
}

const server = new grpc.Server();

server.addService(product_proto.product.DiscountService.service, {discount: discount})
server.bindAsync(LISTEN_ADDR, grpc.ServerCredentials.createInsecure(), () => {
    server.start();
    return this;
  }
);

process.on('SIGTERM', function () {

    console.log("SIGTERM received, shutting down...")
    server.tryShutdown( function() {
        console.log("SIGTERM received, shutting down...done.")
        process.exit(0);
    });

});

module.exports = { checkSameMonthDay, userMonthDayOfBirth, checkUserBirthday, getProduct, findDiscount }