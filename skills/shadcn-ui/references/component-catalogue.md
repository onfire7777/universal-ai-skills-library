# Component Catalogue

The ~15 most-used shadcn/ui components with install commands, key props, and gotchas. Not exhaustive — see [shadcn/ui docs](https://ui.shadcn.com) for the full list.

## Button

```bash
pnpm dlx shadcn@latest add button
```

**Variants**: `default`, `destructive`, `outline`, `secondary`, `ghost`, `link`
**Sizes**: `default`, `sm`, `lg`, `icon`

```tsx
<Button variant="outline" size="sm" onClick={handleClick}>Save</Button>
<Button variant="ghost" size="icon"><Trash className="h-4 w-4" /></Button>
<Button disabled={isPending}>{isPending ? 'Saving...' : 'Save'}</Button>
```

## Input + Label

```bash
pnpm dlx shadcn@latest add input label
```

```tsx
<div className="space-y-2">
  <Label htmlFor="email">Email</Label>
  <Input id="email" type="email" placeholder="you@example.com" />
</div>
```

**Note**: When using with react-hook-form, don't spread `{...field}` — pass props individually to avoid null value issues.

## Card

```bash
pnpm dlx shadcn@latest add card
```

```tsx
<Card>
  <CardHeader>
    <CardTitle>Title</CardTitle>
    <CardDescription>Description</CardDescription>
  </CardHeader>
  <CardContent>Body</CardContent>
  <CardFooter>Footer</CardFooter>
</Card>
```

## Form

```bash
pnpm dlx shadcn@latest add form
pnpm add react-hook-form zod @hookform/resolvers
```

Wraps react-hook-form with shadcn styling. See [recipes.md](recipes.md) for complete form examples.

**Key exports**: `Form`, `FormField`, `FormItem`, `FormLabel`, `FormControl`, `FormMessage`

## Dialog

```bash
pnpm dlx shadcn@latest add dialog
```

```tsx
<Dialog open={open} onOpenChange={setOpen}>
  <DialogTrigger asChild><Button>Open</Button></DialogTrigger>
  <DialogContent className="sm:max-w-md">
    <DialogHeader>
      <DialogTitle>Title</DialogTitle>
      <DialogDescription>Description</DialogDescription>
    </DialogHeader>
    {/* content */}
    <DialogFooter>
      <Button onClick={() => setOpen(false)}>Close</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

**Gotcha**: Override width with `sm:max-w-*` (must match breakpoint prefix).

## Sheet

```bash
pnpm dlx shadcn@latest add sheet
```

Side panel — commonly used for mobile navigation.

```tsx
<Sheet>
  <SheetTrigger asChild><Button variant="ghost" size="icon"><Menu /></Button></SheetTrigger>
  <SheetContent side="left">
    <SheetHeader><SheetTitle>Navigation</SheetTitle></SheetHeader>
    {/* nav links */}
  </SheetContent>
</Sheet>
```

**Sides**: `left`, `right`, `top`, `bottom`

## Table

```bash
pnpm dlx shadcn@latest add table
```

Static table. For sortable/filterable data tables, also install `@tanstack/react-table`. See [recipes.md](recipes.md) for the data table pattern.

```tsx
<Table>
  <TableHeader>
    <TableRow><TableHead>Name</TableHead><TableHead>Email</TableHead></TableRow>
  </TableHeader>
  <TableBody>
    {users.map(u => (
      <TableRow key={u.id}>
        <TableCell>{u.name}</TableCell>
        <TableCell>{u.email}</TableCell>
      </TableRow>
    ))}
  </TableBody>
</Table>
```

## Select

```bash
pnpm dlx shadcn@latest add select
```

```tsx
<Select value={value} onValueChange={setValue}>
  <SelectTrigger><SelectValue placeholder="Choose..." /></SelectTrigger>
  <SelectContent>
    <SelectItem value="option1">Option 1</SelectItem>
    <SelectItem value="option2">Option 2</SelectItem>
  </SelectContent>
</Select>
```

**Gotcha**: No empty string values. Use `"__any__"` sentinel for "All" options.

## Toast (Sonner)

```bash
pnpm dlx shadcn@latest add toast
pnpm add sonner
```

Add `<Toaster />` to your root layout, then:

```tsx
import { toast } from 'sonner'

toast.success('Saved successfully')
toast.error('Something went wrong')
toast.promise(saveData(), {
  loading: 'Saving...',
  success: 'Saved!',
  error: 'Failed to save',
})
```

## Tabs

```bash
pnpm dlx shadcn@latest add tabs
```

```tsx
<Tabs defaultValue="general">
  <TabsList>
    <TabsTrigger value="general">General</TabsTrigger>
    <TabsTrigger value="security">Security</TabsTrigger>
  </TabsList>
  <TabsContent value="general">General settings...</TabsContent>
  <TabsContent value="security">Security settings...</TabsContent>
</Tabs>
```

## Dropdown Menu

```bash
pnpm dlx shadcn@latest add dropdown-menu
```

```tsx
<DropdownMenu>
  <DropdownMenuTrigger asChild><Button variant="ghost" size="icon"><MoreHorizontal /></Button></DropdownMenuTrigger>
  <DropdownMenuContent align="end">
    <DropdownMenuItem onClick={handleEdit}>Edit</DropdownMenuItem>
    <DropdownMenuSeparator />
    <DropdownMenuItem className="text-destructive" onClick={handleDelete}>Delete</DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
```

## Badge

```bash
pnpm dlx shadcn@latest add badge
```

**Variants**: `default`, `secondary`, `outline`, `destructive`

```tsx
<Badge variant="secondary">Draft</Badge>
<Badge variant="destructive">Overdue</Badge>
```

## Switch

```bash
pnpm dlx shadcn@latest add switch
```

```tsx
<div className="flex items-center gap-2">
  <Switch id="notifications" checked={enabled} onCheckedChange={setEnabled} />
  <Label htmlFor="notifications">Enable notifications</Label>
</div>
```

## Separator

```bash
pnpm dlx shadcn@latest add separator
```

```tsx
<Separator />                    {/* horizontal */}
<Separator orientation="vertical" className="h-6" />  {/* vertical */}
```

## Avatar

```bash
pnpm dlx shadcn@latest add avatar
```

```tsx
<Avatar>
  <AvatarImage src={user.avatar} alt={user.name} />
  <AvatarFallback>{user.name[0]}</AvatarFallback>
</Avatar>
```
