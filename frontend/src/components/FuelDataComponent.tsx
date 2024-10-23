// FuelDataComponent.tsx
import React, { useState, useEffect } from 'react';
import { firestore } from "../firebase.ts";
import { collection, query, orderBy, limit, getDocs, where, Timestamp } from 'firebase/firestore';

// Define the structure of your fuel data
interface FuelData {
    time: Timestamp;
    storeID: string;
    ean: string;
    price: number;
    priceDate: Timestamp;
    state: string;
    suburb: string;
    address: string;
    postcode: string;
}

// Map EAN codes to their descriptive names
const eans: { [key: string]: string } = {
    "52": "Special Unleaded 91",
    "53": "Special Diesel",
    "57": "Special E10",
    "56": "Supreme+ 98",
    "55": "Extra 95",
    "54": "LPG",
    // Add other EANs as needed
};

const FuelDataComponent: React.FC = () => {
    // State to hold the cheapest data: state -> ean -> {price, fuelDataList}
    const [stateEanCheapestMap, setStateEanCheapestMap] = useState<{
        [state: string]: {
            [ean: string]: {
                price: number;
                fuelDataList: FuelData[];
            };
        };
    }>({});

    useEffect(() => {
        const fetchData = async () => {
            try {
                // Step 1: Get the latest 'time' value
                const fuelCollection = collection(firestore, 'fuel');
                const latestTimeQuery = query(fuelCollection, orderBy('time', 'desc'), limit(1));
                const latestTimeSnapshot = await getDocs(latestTimeQuery);

                if (latestTimeSnapshot.empty) {
                    console.log('No data found');
                    return;
                }

                const latestTimeDoc = latestTimeSnapshot.docs[0];
                const latestTime = latestTimeDoc.data().time as Timestamp;

                // Step 2: Get all documents with that 'time' value
                const dataQuery = query(fuelCollection, where('time', '==', latestTime));
                const dataSnapshot = await getDocs(dataQuery);

                const fuelDataArray: FuelData[] = [];
                dataSnapshot.forEach((doc) => {
                    fuelDataArray.push(doc.data() as FuelData);
                });

                // Step 3: Organize data into state -> ean -> {price, fuelDataList}
                const tempStateEanCheapestMap: {
                    [state: string]: {
                        [ean: string]: {
                            price: number;
                            fuelDataList: FuelData[];
                        };
                    };
                } = {};

                fuelDataArray.forEach((fuelData) => {
                    const state = fuelData.state;
                    const ean = fuelData.ean;

                    if (!tempStateEanCheapestMap[state]) {
                        tempStateEanCheapestMap[state] = {};
                    }

                    if (!tempStateEanCheapestMap[state][ean]) {
                        tempStateEanCheapestMap[state][ean] = {
                            price: fuelData.price,
                            fuelDataList: [fuelData],
                        };
                    } else {
                        const existingEntry = tempStateEanCheapestMap[state][ean];
                        if (fuelData.price < existingEntry.price) {
                            // Replace with new cheaper price and reset fuelDataList
                            existingEntry.price = fuelData.price;
                            existingEntry.fuelDataList = [fuelData];
                        } else if (fuelData.price === existingEntry.price) {
                            // Same price, add to fuelDataList
                            existingEntry.fuelDataList.push(fuelData);
                        }
                    }
                });

                setStateEanCheapestMap(tempStateEanCheapestMap);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, []);

    return (
        <div>
            <h1>Fuel Prices</h1>
            {Object.keys(stateEanCheapestMap).sort().map((state) => {
                const eanMap = stateEanCheapestMap[state];
                const eansList = Object.keys(eanMap).sort((a, b) => {
                    const priceA = eanMap[a].price;
                    const priceB = eanMap[b].price;
                    return priceA - priceB;
                });

                return (
                    <div key={state} style={{ marginBottom: '40px' }}>
                        <h2>Best prices ({state}):</h2>
                        <table
                            border={1}
                            cellPadding={5}
                            cellSpacing="0"
                            style={{ width: '100%', borderCollapse: 'collapse' }}
                        >
                            <thead>
                            <tr>
                                <th>Fuel</th>
                                <th>Price</th>
                                <th>Address</th>
                                <th>Suburb</th>
                                <th>State</th>
                                <th>Postcode</th>
                                <th>Time scraped</th>
                                <th>Price date</th>
                            </tr>
                            </thead>
                            <tbody>
                            {eansList.map((ean) => {
                                const { price, fuelDataList } = eanMap[ean];

                                // Display all stations with the cheapest price for this EAN
                                return fuelDataList.map((fuelData, index) => (
                                    <tr key={`${state}-${ean}-${index}`}>
                                        <td>{eans[ean] || ean}</td>
                                        <td>{(price / 10).toFixed(2)}</td>
                                        <td>{fuelData.address}</td>
                                        <td>{fuelData.suburb}</td>
                                        <td>{fuelData.state}</td>
                                        <td>{fuelData.postcode}</td>
                                        <td>{fuelData.time.toDate().toLocaleString()}</td>
                                        <td>{fuelData.priceDate.toDate().toLocaleString()}</td>
                                    </tr>
                                ));
                            })}
                            </tbody>
                        </table>
                    </div>
                );
            })}
        </div>
    );
};

export default FuelDataComponent;
