db.getCollection("customers").find(
    {'_id': {$eq: 'ryann.graham'}},
    {'customer_id': 1, 'name': 1}
)

db.getCollection("cases").find(
    {'customer_id': {$eq: 'dawson.maggio'}},
    {'case_id': 1, 'customer_id': 1}
)

db.getCollection("customers").aggregate([
    {$match: {'_id': {$eq: 'dawson.maggio'}}},
    {
        $lookup: {
            from: 'cases',
            localField: '_id',
            foreignField: 'customer_id',
            as: 'cases'
        }
    },
    {
        $lookup: {
            from: 'accounts',
            localField: '_id',
            foreignField: 'customer_id',
            as: 'accounts'
        }
    }
])