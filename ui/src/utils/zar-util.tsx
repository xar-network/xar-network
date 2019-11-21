const ZarClient = require('@zar-network/javascript-sdk');

// const zarApi = 'https://smhuzpepz8.execute-api.eu-west-1.amazonaws.com/testnet'
const zarApi = 'http://localhost:1317'

export const unlockAccount = async (params:any) => {
  try {
    const client = new ZarClient(zarApi)
    await client.initChain()

    const acc = await client.recoverAccountFromKeystore(params.keystore, params.password);
    return acc
  } catch (err) {
    throw err;
  }
};

export const createAccount = async (params:any) => {
  try {
    const client = new ZarClient(zarApi)
    await client.initChain()

    const acc = await client.createAccountWithMneomnic()
    //add error checking

    const keystore = ZarClient.crypto.generateKeyStore(acc.privateKey, params.password)
    acc.keystore = keystore

    return acc

  } catch (err) {
    throw err;
  }
};
