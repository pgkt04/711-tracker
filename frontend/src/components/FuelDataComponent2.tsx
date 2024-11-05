import React, { useState, useEffect } from 'react';
import { firestore } from "../firebase.ts";
import { collection, query, orderBy, limit, getDocs, where, Timestamp } from 'firebase/firestore';
import { ChevronRight } from "lucide-react";

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

const eans: { [key: string]: string } = {
    "52": "Special Unleaded 91",
    "53": "Special Diesel",
    "57": "Special E10",
    "56": "Supreme+ 98",
    "55": "Extra 95",
    "54": "LPG",
};

const FuelDataComponent: React.FC = () => {
    const [stateEanCheapestMap, setStateEanCheapestMap] = useState<{
        [state: string]: {
            [ean: string]: {
                price: number;
                fuelDataList: FuelData[];
            };
        };
    }>({});

    // Single page state for each state
    const [statePages, setStatePages] = useState<{
        [state: string]: number;
    }>({});

    useEffect(() => {
        const fetchData = async () => {
            try {
                const fuelCollection = collection(firestore, 'fuel');
                const latestTimeQuery = query(fuelCollection, orderBy('time', 'desc'), limit(1));
                const latestTimeSnapshot = await getDocs(latestTimeQuery);

                if (latestTimeSnapshot.empty) {
                    console.log('No data found');
                    return;
                }

                const latestTimeDoc = latestTimeSnapshot.docs[0];
                const latestTime = latestTimeDoc.data().time as Timestamp;

                const dataQuery = query(fuelCollection, where('time', '==', latestTime));
                const dataSnapshot = await getDocs(dataQuery);

                const fuelDataArray: FuelData[] = [];
                dataSnapshot.forEach((doc) => {
                    fuelDataArray.push(doc.data() as FuelData);
                });

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
                            existingEntry.price = fuelData.price;
                            existingEntry.fuelDataList = [fuelData];
                        } else if (fuelData.price === existingEntry.price) {
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

    // Check if a state has any duplicates
    const hasStateDuplicates = (state: string): boolean => {
        const eanMap = stateEanCheapestMap[state];
        return Object.values(eanMap).some(({ fuelDataList }) => fuelDataList.length > 1);
    };

    // Toggle between pages for a state
    const toggleStatePage = (state: string) => {
        setStatePages(prev => ({
            ...prev,
            [state]: (prev[state] || 0) === 0 ? 1 : 0
        }));
    };

    return (
        <div className="max-w-6xl mx-auto p-4">
            <h1 className="text-2xl font-bold mb-6">7/11 Fuel Prices</h1>

            {Object.keys(stateEanCheapestMap).sort().map((state) => {
                const eanMap = stateEanCheapestMap[state];
                const eansList = Object.keys(eanMap).sort((a, b) => {
                    const priceA = eanMap[a].price;
                    const priceB = eanMap[b].price;
                    return priceA - priceB;
                });
                const currentPage = statePages[state] || 0;
                const showDuplicates = currentPage === 1;

                return (
                    <div key={state} className="mb-8">
                        <h2 className="text-xl font-semibold mb-4">Best prices ({state}):</h2>
                        <div className="overflow-x-auto">
                            <table className="min-w-full border-collapse border border-gray-200">
                                <thead>
                                    <tr className="bg-gray-50">
                                        <th className="border p-2 text-left">Fuel</th>
                                        <th className="border p-2 text-left">Price</th>
                                        <th className="border p-2 text-left">Address</th>
                                        <th className="border p-2 text-left">Suburb</th>
                                        <th className="border p-2 text-left">State</th>
                                        <th className="border p-2 text-left">Postcode</th>
                                        <th className="border p-2 text-left">Time scraped</th>
                                        <th className="border p-2 text-left">Price date</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {eansList.map((ean) => {
                                        const { price, fuelDataList } = eanMap[ean];

                                        if (showDuplicates) {
                                            // Show only entries with duplicates, starting from the second one
                                            if (fuelDataList.length <= 1) return null;
                                            return fuelDataList.slice(1).map((fuelData, idx) => (
                                                <tr key={`${state}-${ean}-${idx + 1}`}>
                                                    <td className="border p-2">{eans[ean] || ean}</td>
                                                    <td className="border p-2">{(price / 10).toFixed(2)}</td>
                                                    <td className="border p-2">{fuelData.address}</td>
                                                    <td className="border p-2">{fuelData.suburb}</td>
                                                    <td className="border p-2">{fuelData.state}</td>
                                                    <td className="border p-2">{fuelData.postcode}</td>
                                                    <td className="border p-2">
                                                        {fuelData.time.toDate().toLocaleString()}
                                                    </td>
                                                    <td className="border p-2">
                                                        {fuelData.priceDate.toDate().toLocaleString()}
                                                    </td>
                                                </tr>
                                            ));
                                        } else {
                                            // Show only the first entry for each fuel type
                                            const fuelData = fuelDataList[0];
                                            return (
                                                <tr key={`${state}-${ean}`}>
                                                    <td className="border p-2">{eans[ean] || ean}</td>
                                                    <td className="border p-2">{(price / 10).toFixed(2)}</td>
                                                    <td className="border p-2">{fuelData.address}</td>
                                                    <td className="border p-2">{fuelData.suburb}</td>
                                                    <td className="border p-2">{fuelData.state}</td>
                                                    <td className="border p-2">{fuelData.postcode}</td>
                                                    <td className="border p-2">
                                                        {fuelData.time.toDate().toLocaleString()}
                                                    </td>
                                                    <td className="border p-2">
                                                        {fuelData.priceDate.toDate().toLocaleString()}
                                                    </td>
                                                </tr>
                                            );
                                        }
                                    })}
                                </tbody>
                            </table>
                        </div>

                        {/* Single next button that only appears if there are duplicates */}
                        {hasStateDuplicates(state) && (
                            <button
                                onClick={() => toggleStatePage(state)}
                                className="flex items-center gap-2 mt-4 px-4 py-2 bg-blue-50 hover:bg-blue-100 rounded-md text-blue-600"
                            >
                                {showDuplicates ? (
                                    "Show primary locations"
                                ) : (
                                    <>
                                        Show additional locations
                                        <ChevronRight className="w-4 h-4" />
                                    </>
                                )}
                            </button>
                        )}
                    </div>
                );
            })}
        </div>
    );
};

export default FuelDataComponent;