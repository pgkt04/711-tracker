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
    const [stateEanCheapestMap, setStateEanCheapestMap] = useState<{
        [state: string]: {
            [index: number]: {
                [ean: string]: {
                    price: number;
                    fuelDataList: FuelData[];
                };
            };
        };
    }>({});

    // New state to keep track of current index per state
    const [currentIdxMap, setCurrentIdxMap] = useState<{ [state: string]: number }>({});

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
                        [idx: number]: {
                            [ean: string]: {
                                price: number;
                                fuelDataList: FuelData[];
                            };
                        };
                    };
                } = {};

                // Sort fuelDataArray by price
                fuelDataArray.sort((a, b) => a.price - b.price);

                // Insert into tempStateEanCheapestMap
                fuelDataArray.forEach((fuelData) => {
                    const state = fuelData.state;
                    const ean = fuelData.ean;

                    if (!tempStateEanCheapestMap[state]) {
                        tempStateEanCheapestMap[state] = {};
                    }

                    let inserted = false;
                    let idx = 0;
                    while (!inserted) {
                        if (!tempStateEanCheapestMap[state][idx]) {
                            tempStateEanCheapestMap[state][idx] = {}
                        }

                        if (!tempStateEanCheapestMap[state][idx][ean]) {
                            tempStateEanCheapestMap[state][idx][ean] = {
                                price: fuelData.price,
                                fuelDataList: [fuelData],
                            };
                            inserted = true;
                            break;
                        }
                        idx++;
                    }
                });

                // **Step 4: Sort EANs Lexicographically within Each State and Index**
                for (const state in tempStateEanCheapestMap) {
                    for (const idx in tempStateEanCheapestMap[state]) {
                        // Get the EANs and sort them lexicographically based on their descriptive names
                        const sortedEans = Object.keys(tempStateEanCheapestMap[state][idx]).sort((a, b) => {
                            const nameA = eans[a] || a;
                            const nameB = eans[b] || b;
                            return nameA.localeCompare(nameB);
                        });

                        // Create a new sorted object
                        const sortedEanObj: {
                            [ean: string]: {
                                price: number;
                                fuelDataList: FuelData[];
                            };
                        } = {};

                        sortedEans.forEach((ean) => {
                            sortedEanObj[ean] = tempStateEanCheapestMap[state][idx][ean];
                        });

                        // Replace the unsorted object with the sorted one
                        tempStateEanCheapestMap[state][idx] = sortedEanObj;
                    }
                }

                setStateEanCheapestMap(tempStateEanCheapestMap);

                // Initialize currentIdxMap with 0 for each state
                const initialIdxMap: { [state: string]: number } = {};
                Object.keys(tempStateEanCheapestMap).forEach(state => {
                    initialIdxMap[state] = 0;
                });

                setCurrentIdxMap(initialIdxMap);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, []);

    // Handler for the Next button
    const handleNext = (state: string) => {
        setCurrentIdxMap(prev => {
            const currentIdx = prev[state] || 0;
            const maxIdx = Object.keys(stateEanCheapestMap[state] || {}).length - 1;
            const nextIdx = currentIdx >= maxIdx ? 0 : currentIdx + 1;
            return { ...prev, [state]: nextIdx };
        });
    };

    return (
        <div className="container mx-auto p-6">
            <h1 className="text-3xl font-bold mb-8 text-center text-blue-600">7/11 Fuel Prices</h1>
            {Object.keys(stateEanCheapestMap).sort().map((state) => {
                const eanMapForState = stateEanCheapestMap[state];
                const currentIdx = currentIdxMap[state] || 0;
                const eanMap = eanMapForState[currentIdx] || {};
                const eansList = Object.keys(eanMap).sort((a, b) => {
                    const nameA = eans[a] || a;
                    const nameB = eans[b] || b;
                    return nameA.localeCompare(nameB);
                });

                // Calculate the total number of idx for this state
                const totalIdx = Object.keys(eanMapForState).length;

                return (
                    <div key={state} className="mb-12">
                        <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                            Best Prices in <span className="text-blue-500">{state}</span> - Page {currentIdx + 1} of {totalIdx}
                        </h2>
                        <div className="min-w-full bg-white rounded-lg shadow overflow-hidden">
                            <table className="min-w-full">
                                <thead className="bg-gray-200">
                                    <tr>
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Fuel</th>
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Price ($)</th>
                                        {/* <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Address</th> */}
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Suburb</th>
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">State</th>
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Postcode</th>
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Time Scraped</th>
                                        <th className="py-3 px-5 text-left text-sm font-medium text-gray-700 border-b border-gray-300">Price Date</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-300">
                                    {eansList.map((ean) => {
                                        const { price, fuelDataList } = eanMap[ean];

                                        // Display all stations with the cheapest price for this EAN
                                        return fuelDataList.map((fuelData, index, array) => {
                                            // Determine if this is the last row in the entire table
                                            // Calculate the total number of rows for this state
                                            const totalRows = eansList.reduce((acc, currentEan) => acc + (stateEanCheapestMap[state][currentIdx][currentEan].fuelDataList.length), 0);
                                            // Calculate the current row index
                                            let currentRowIndex = 0;
                                            for (let e of eansList) {
                                                if (e === ean) break;
                                                currentRowIndex += stateEanCheapestMap[state][currentIdx][e].fuelDataList.length;
                                            }
                                            currentRowIndex += index;
                                            const isLastRow = currentRowIndex === totalRows - 1;

                                            return (
                                                <tr
                                                    key={`${state}-${ean}-${index}`}
                                                    className={`${index % 2 === 0 ? 'bg-white' : 'bg-gray-50'} border-b border-gray-200 ${isLastRow ? 'border-b-0' : ''}`}
                                                >
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{eans[ean] || ean}</td>
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">${(price / 10).toFixed(2)}</td>
                                                    {/* <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{fuelData.address}</td> */}
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{fuelData.suburb}</td>
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{fuelData.state}</td>
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{fuelData.postcode}</td>
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{fuelData.time.toDate().toLocaleString()}</td>
                                                    <td className="py-4 px-5 text-sm text-gray-700 border-gray-200">{fuelData.priceDate.toDate().toLocaleString()}</td>
                                                </tr>
                                            );
                                        });
                                    })}
                                </tbody>
                            </table>
                        </div>
                        {totalIdx > 1 && (
                            <div className="flex justify-end mt-4">
                                <button
                                    onClick={() => handleNext(state)}
                                    className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                                >
                                    Next
                                </button>
                            </div>
                        )}
                    </div>
                );
            })}
        </div>
    );
};

export default FuelDataComponent;
