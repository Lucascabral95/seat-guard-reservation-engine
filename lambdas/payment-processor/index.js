exports.handler = async (event) => {
  for (const record of event.Records) {
    const body = JSON.parse(record.body); 
    
    console.log("Message body:", body);
  }

  return;
};
