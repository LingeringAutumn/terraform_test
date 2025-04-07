exports.handler = async (event) => {
    console.log('事件对象:', JSON.stringify(event, null, 2));
    
    const response = {
        statusCode: 200,
        body: JSON.stringify({
            message: 'Hello, World! Lambda函数已成功执行。',
            timestamp: new Date().toISOString(),
            event: event
        }),
    };
    
    return response;
};