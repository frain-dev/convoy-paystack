require 'bundler/setup'
require 'paystack'

module API
  INVOICE_PATH = "/paymentrequest"
end

class PaystackInvoice < PaystackBaseObject
  def create(data={})
    return PaystackInvoice.create(@paystack, data)
  end

  def self.create(paystackObj, data)
		initPostRequest(paystackObj, "#{API::INVOICE_PATH}", data, true)
	end
end