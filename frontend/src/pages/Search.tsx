import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Expert, IscedField, getExperts, getIscedFields } from "../api/api";

export default function Search() {
  const [experts, setExperts] = useState<Expert[]>([]);
  const [iscedFields, setIscedFields] = useState<IscedField[]>([]);
  const [filters, setFilters] = useState({
    name: "",
    affiliation: "",
    is_bahraini: undefined as boolean | undefined,
    isced_field_id: "",
    is_available: true,
    page: 1,
    limit: 10,
  });
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [expertsData, fieldsData] = await Promise.all([
          getExperts(filters),
          getIscedFields(),
        ]);
        setExperts(expertsData);
        setIscedFields(fieldsData);
      } catch {
        setError("Failed to load data");
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [filters]);

  const handleSort = () => {
    const newOrder = sortOrder === "asc" ? "desc" : "asc";
    setSortOrder(newOrder);
    setExperts([...experts].sort((a, b) =>
      newOrder === "asc"
        ? a.name.localeCompare(b.name)
        : b.name.localeCompare(a.name)
    ));
  };

  return (
    <div className="p-6">
      <h2 className="text-2xl mb-4">Search Experts</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <Input
          placeholder="Search by name..."
          value={filters.name}
          onChange={(e) => setFilters({ ...filters, name: e.target.value, page: 1 })}
        />
        <Input
          placeholder="Filter by affiliation..."
          value={filters.affiliation}
          onChange={(e) => setFilters({ ...filters, affiliation: e.target.value, page: 1 })}
        />
        <Select
          value={filters.isced_field_id}
          onValueChange={(value) => setFilters({ ...filters, isced_field_id: value, page: 1 })}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select ISCED Field" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">All Fields</SelectItem>
            {iscedFields.map((field) => (
              <SelectItem key={field.id} value={field.id}>
                {field.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select
          value={filters.is_bahraini === undefined ? "" : filters.is_bahraini ? "true" : "false"}
          onValueChange={(value) => setFilters({
            ...filters,
            is_bahraini: value === "" ? undefined : value === "true",
            page: 1,
          })}
        >
          <SelectTrigger>
            <SelectValue placeholder="Bahraini Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">All</SelectItem>
            <SelectItem value="true">Bahraini</SelectItem>
            <SelectItem value="false">Non-Bahraini</SelectItem>
          </SelectContent>
        </Select>
        <Select
          value={filters.is_available ? "true" : "false"}
          onValueChange={(value) => setFilters({
            ...filters,
            is_available: value === "true",
            page: 1,
          })}
        >
          <SelectTrigger>
            <SelectValue placeholder="Availability" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="true">Available</SelectItem>
            <SelectItem value="false">Not Available</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {loading && <p>Loading...</p>}
      {error && <p className="text-red-500">{error}</p>}
      {!loading && !error && (
        <>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead onClick={handleSort} className="cursor-pointer">
                  Name {sortOrder === "asc" ? "↑" : "↓"}
                </TableHead>
                <TableHead>Affiliation</TableHead>
                <TableHead>Bahraini</TableHead>
                <TableHead>ISCED Field</TableHead>
                <TableHead>Available</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {experts.map((expert) => (
                <TableRow key={expert.id}>
                  <TableCell>{expert.name}</TableCell>
                  <TableCell>{expert.affiliation}</TableCell>
                  <TableCell>{expert.is_bahraini ? "Yes" : "No"}</TableCell>
                  <TableCell>{expert.isced_field_id}</TableCell>
                  <TableCell>{expert.is_available ? "Yes" : "No"}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <div className="flex justify-between mt-4">
            <Button
              disabled={filters.page === 1}
              onClick={() => setFilters({ ...filters, page: filters.page - 1 })}
            >
              Previous
            </Button>
            <span>Page {filters.page}</span>
            <Button
              disabled={experts.length < filters.limit}
              onClick={() => setFilters({ ...filters, page: filters.page + 1 })}
            >
              Next
            </Button>
          </div>
        </>
      )}
    </div>
  );
}