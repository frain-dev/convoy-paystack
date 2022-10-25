require_relative './invoices.rb'

payload = {
  "line_items": [
    {"name": "item 1", "amount": 200000},
    {"name": "item 2", "amount": 20000}
  ],
  "due_date": "2022-12-08"
}

customers = [
  {
    "name": "Subomi",
    "customer_id": "CUS_qaqds0b4n9597a2"
  },
  {
    "name": "Emmanuel",
    "customer_id": "CUS_vz25k1yvtaeumgp"
  },
  {
    "name": "Raymond",
    "customer_id": "CUS_l8vlbjnrk6gkqgw"
  },
  {
    "name": "Andrew",
    "customer_id": "CUS_972803vi04si3gu"
  }
]

paystack_obj = Paystack.new(ENV['PUBLIC_KEY'], ENV['SECRET_KEY'])
invoice_count = ENV['INVOICE_COUNT']

invoice_count.times do 
  idx = rand(0...4)
  payload["description"] = customers[idx][:name] + "'s invoice"
  payload["customer"] = customers[idx][:customer_id]

  invoice = PaystackInvoice.create(paystack_obj, payload)
  sleep 1
end
