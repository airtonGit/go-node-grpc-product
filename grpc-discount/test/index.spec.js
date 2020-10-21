const chai = require('chai');
const chaiSubset = require('chai-subset');
const subSet = require('chai-subset'); 

chai.use(subSet);

const index = require('../index'); // Arquivo a ser testado

const productSchema = {
    id: id => id,
    price_in_cents: price_in_cents => price_in_cents
};

describe('Functions tests', () => {

    it('checkSameMonthDay', () => {

        let today = new Date();
        let got = index.checkSameMonthDay(`${today.getMonth()+1}/${today.getDate()}`, today)
        
        chai.expect(got).to.be.true
        
    });

    it('userMonthDayOfBirth', () => {
        got = index.userMonthDayOfBirth("11/27/1983")

        chai.expect(got).to.be.equal("11/27")
    })

    it('checkUserBirthday', () => {
        let birthday = new Date();
        birthday.setMonth(9)
        birthday.setDate(20)

        let got = index.checkUserBirthday("3", birthday)
        
        chai.expect(got).to.be.true
    })

    it('getProduct', () => {
        got = index.getProduct("1")

        chai.expect(got).to.be.not.undefined
        chai.expect(got).to.containSubset(productSchema);
        chai.expect(got.id).to.string("1")
    })

    it('findDiscount', () => {
        got = index.findDiscount('1', '1');

        chai.expect(got).to.be.not.undefined
        chai.expect(got).to.containSubset(productSchema);
        chai.expect(got.id).to.string("1")
    })

});