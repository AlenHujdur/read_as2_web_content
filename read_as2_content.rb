require 'json'
require 'csv'
require 'digest/sha1'
require "base64"
require 'charlock_holmes'
require 'tempfile'
require "zlib"
require "prawn"
require 'stringio'
require 'date'

public def is_i?
  !!(self =~ /\A[-+]?[0-9]+\z/)
end
def getKeys(folder = "db", type)
  file = File.read("../../../../../../../Desktop/#{folder}/ImportIDs.txt")
  data_hash = JSON.parse(file)
  if type == "betrokkenen"
    CSV.open("#{folder}_betrokkenen_keys.csv", "w") do |csv|
      csv << ["brio_bet_id;sigura_bet_id"]
      data_hash.each do |key, value|
        if key == "betrokkenen"
          value.each do |betrokken|
            csv << ["#{betrokken[0]}; #{betrokken[1]}"]
          end
        end
      end
    end
  end
  if type == "contracten"
    CSV.open("#{folder}_contracten_keys.csv", "w") do |csv|
      csv << ["brio_contract_id;sigura_contract_id"]
      data_hash.each do |key, value|
        if key == "contracten"
          value.each do |contracten|
            csv << ["#{contracten[0]}; #{contracten[1]}"]
          end
        end
      end
    end
  end
  if type == "schadegevallen"
    CSV.open("#{folder}_schadegevallen_keys.csv", "w") do |csv|
      csv << ["brio_schade_id;sigura_schade_id"]
      data_hash.each do |key, value|
        if key == "schadegeval"
          value.each do |schadegeval|
            csv << ["#{schadegeval[0]}; #{schadegeval[1]}"]
          end
        end
      end
    end
  end
  if type == "risicos"
    CSV.open("#{folder}_risicos_keys.csv", "w") do |csv|
      csv << ["brio_risico_id;sigura_risico_id"]
      data_hash.each do |key, value|
        if key == "risicos"
          value.each do |risico|
            csv << ["#{risico[0]}; #{risico[1]}"]
          end
        end
      end
    end
  end
  if type == "bestanden"
    CSV.open("#{folder}_bestanden_keys.csv", "w") do |csv|
      csv << ["brio_bestand_id;sigura_bestand_id"]
      data_hash.each do |key, value|
        if key == "bestanden"
          value.each do |bestand|
            csv << ["#{bestand[0]}; #{bestand[1]}"]
          end
        end
      end
    end
  end
  if type == "addressen"
    CSV.open("#{folder}_addressen_keys.csv", "w") do |csv|
      csv << ["brio_address_id;sigura_address_id"]
      data_hash.each do |key, value|
        if key == "addresen"
          value.each do |bestand|
            csv << ["#{bestand[0]}; #{bestand[1]}"]
          end
        end
      end
    end
  end
end
def read_coulmns(filename, start, kill)
  storage = []
  open(filename) do |f|
    state = nil

    while (line = f.gets)
      case (state)
        when nil,' '
          # Look for the line beginning with "FILE_SOURCES"
          if (line.match(/^#{start}/))
            state = :sources
            #storage << line.strip!.delete_prefix('"').delete_suffix('"')
          end
        when :sources
          # Stop printing if you hit something starting with "END"
          if (line.match(/^#{kill}/))
            state = nil
          elsif line.strip != ""
              storage << line.strip!.delete_prefix('"').delete_suffix('"')
          end
        end
      end
    f.close
  end

  storage.each do |line|
    puts line
  end
  puts storage.count
end
def read_table(filename = 'msg_as2_att2.txt',endline)
  f = File.open(filename, 'r', :encoding => 'utf-8')
  content = []
  files = []
  f.each_line do |line|
    line = line.encode("UTF-8", :invalid => :replace, :undef => :replace, :replace => "?")
    #d = CharlockHolmes::EncodingDetector.detect(line)
    #line = line.encode("UTF-8", d[:encoding], invalid: :replace, replace: "")
    content << [line.strip.split(";,")] unless !!(line =~ /^#{endline}/)
  end
  #content = content[0..1]
  content.each do |line|
    line.each do |element|
      files << element
    end
  end
  f.close
  i = 0
  #puts files[15101]
  #exit
  #2335
  #puts files.inspect
  #exit

  tmp_filename = ''
  tmp_content = ''
  tmp_file_id = ''
  tmp_type= ''
  i = 0
  if filename.include? 'address.txt'
    folder = "schaekers_brio"
    #DB60324
    f = File.open("../../../../../../../Desktop/#{folder}/location.txt", 'r', :encoding => 'utf-8')
    content = []
    locations = []
    f.each_line do |line|
      line = line.encode("UTF-8", :invalid => :replace, :undef => :replace, :replace => "?")
      content << [line.strip.split(";,")] unless !!(line =~ /^#{endline}/)
    end
    content.each do |line|
      line.each do |element|
        locations << element
      end
    end
    f.close

    CSV.open("addresses.csv", "w") do |csv|
      csv << ["p_address;city;street;house_no;c_country;zip;type;status"]
      locations.each do |location|
        files.each do |line|
          #line[0].map {|x| x[/\d+/]}
          #puts "checking: #{line[0]}"
          if line[0] != nil
            #puts "entering: #{line[0]}"
            if line[0].strip.is_i? == true
              if location[2].strip.to_i == line[0].strip.to_i
                puts "Location: #{location[0]}"
                csv << ["#{line[0].strip.to_i};#{line[1].strip};#{line[2].strip};#{line[3]};#{line[5]};#{line[6]};#{line[15].strip};#{location[3].strip}"]
              end
            end
          end
        end
      end
    end
    # locations.each do |line|
    #  puts "#{line[0]} - #{line[1]} - #{line[2]} - #{line[3]}"
    # end


  end
  if filename.include? 'doc_blob_desc.txt'
    CSV.open("doc_blob_desc.csv", "w") do |csv|
      csv << ["p_doc_blob;document_name;file_extension;c_object;c_doc_location;c_doc_type;create_date"]
      files.each do |line|
        #line[0].map {|x| x[/\d+/]}
        date_temp = DateTime.parse(line[15].strip)
        date = date_temp.strftime('%Y-%m-%d %H:%M:%S')
        puts "checking: #{line[0]}"
        if line[0] != nil
          puts "entering: #{line[0]}"
          if line[0].strip.is_i? == true
            puts "number: #{line[0]}"
            csv << ["#{line[0].strip.to_i};#{line[1].strip};#{line[2].strip};#{line[3].strip};#{line[5].strip};#{line[6].strip};#{date}"]
          end
        end
      end
    end
  end
  if filename.include? 'bac.txt'
    CSV.open("bac_files_dates.csv", "w") do |csv|
      csv << ["p_doc_blob;date_blob;create_date;last_date;date_archive"]

      files.each do |line|
        #line[0].map {|x| x[/\d+/]}
        puts "checking: #{line[0]}"
          if line[0] != nil
            puts "entering: #{line[0]}"
            if line[0].strip.is_i? == true
              puts "number: #{line[0]}"
            csv << ["#{line[0]};#{line[2]};#{line[18]};#{line[16]};#{line[25]};"]
            end
          end
      end
  end
  end
  if filename == 'msg_as2_att2.txt'
    Dir.mkdir("sha1files") unless File.exists?("sha1files")
    CSV.open("as2_files.csv", "w") do |csv|
      csv << ["p_msg_as2;type;filename;filesize;sha1hash"]

      files.each do |file|

        if file[5] != nil && file[8].rstrip.length > 76
          File.open("as2_msg_file_#{i}_#{file[5].rstrip}", 'wb') {|f| f.write(Base64.decode64(file[8].rstrip))}
          File.open("sha1files/#{Digest::SHA1.hexdigest(file[8].rstrip)}", 'wb') {|f| f.write(Base64.decode64(file[8].rstrip))}
          csv << ["#{file[1]};#{file[4].rstrip};#{"as2_msg_file_#{i}_#{file[5].rstrip}"}; #{File.size("as2_msg_file_#{i}_#{file[5].rstrip}")}; #{Digest::SHA1.hexdigest(file[8].rstrip)}"]
          tmp_filename = "as2_msg_file_#{i}_#{file[5].rstrip}"
          tmp_file_id = file[1].rstrip
          tmp_type = file[4].rstrip
          i += 1
        else if file[5] != nil && file[8].rstrip.length == 76 && file[1].rstrip != tmp_file_id #&& tmp_file_id != ''
          File.open(tmp_filename, 'wb') {|f| f.write(Base64.decode64(tmp_content))}
          File.open("sha1files/#{Digest::SHA1.hexdigest(tmp_content)}", 'wb') {|f| f.write(Base64.decode64(tmp_content))}
          csv << ["#{tmp_file_id};#{tmp_type};#{tmp_filename}; #{File.size(tmp_filename)}; #{Digest::SHA1.hexdigest(tmp_content)}"]
          tmp_filename = "as2_msg_file_#{i}_#{file[5].rstrip}"
          tmp_content = file[8].rstrip
          tmp_file_id = file[1].rstrip
          tmp_type = file[4].rstrip
          i += 1
        else if file[0].rstrip.length <= 76
          tmp_content += file[0].rstrip
          puts "writing content for #{tmp_filename}"
        end
        end
        end
      end

    end
  end


end
def read_objects(filename = 'msg_as2_att.txt',endline)
  f = File.open(filename, 'r', :encoding => 'utf-8')
  content = []
  files = []
  f.each_line do |line|
    line = line.encode("UTF-8", :invalid => :replace, :undef => :replace, :replace => "?")
    #d = CharlockHolmes::EncodingDetector.detect(line)
    #line = line.encode("UTF-8", d[:encoding], invalid: :replace, replace: "")
    content << [line.strip.split(";,")] unless !!(line =~ /^#{endline}/)
  end
  #content = content[0..1]
  content.each do |line|
    line.each do |element|
      files << element
    end
  end
  f.close
  i = 0

  tmp_filename = ''
  tmp_content = ''
  tmp_file_id = ''
  tmp_type= ''
  i = 0
  Dir.mkdir("sha1files") unless Dir.exist?("sha1files")
  Dir.mkdir("files") unless Dir.exist?("files")
  CSV.open("as2_files.csv", "w") do |csv|
    csv << ["p_msg_as2;type;filename;filesize;sha1hash"]

    files.each do |file|
      if file[5] != nil && file[8].rstrip.length > 76
        File.open("files/as2_msg_file_#{i}_#{file[5].rstrip}", 'wb') {|f| f.write(Base64.decode64(file[8].rstrip))}
        File.open("sha1files/#{Digest::SHA1.hexdigest(file[8].rstrip)}", 'wb') {|f| f.write(Base64.decode64(file[8].rstrip))}
        csv << ["#{file[1]};#{file[4].rstrip};#{"files/as2_msg_file_#{i}_#{file[5].rstrip}"}; #{File.size("files/as2_msg_file_#{i}_#{file[5].rstrip}")}; #{Digest::SHA1.hexdigest(file[8].rstrip)}"]
        tmp_filename = "files/as2_msg_file_#{i}_#{file[5].rstrip}"
        tmp_file_id = file[1].rstrip
        tmp_type = file[4].rstrip
        i += 1
      else if file[5] != nil && file[8].rstrip.length == 76 && file[1].rstrip != tmp_file_id #&& tmp_file_id != ''
             File.open(tmp_filename, 'wb') {|f| f.write(Base64.decode64(tmp_content))}
             File.open("sha1files/#{Digest::SHA1.hexdigest(tmp_content)}", 'wb') {|f| f.write(Base64.decode64(tmp_content))}
             csv << ["#{tmp_file_id};#{tmp_type};#{tmp_filename}; #{File.size(tmp_filename)}; #{Digest::SHA1.hexdigest(tmp_content)}"]
             tmp_filename = "files/as2_msg_file_#{i}_#{file[5].rstrip}"
             tmp_content = file[8].rstrip
             tmp_file_id = file[1].rstrip
             tmp_type = file[4].rstrip
             i += 1

           else if file[0].rstrip.length <= 76
                  tmp_content += file[0].rstrip
                  puts "writing content for #{tmp_filename}"
                end
           end
      end
    end

  end
  puts "done"
end


#filename = "DB60324_ListColumns.out" #'msg_as2_att.txt'
endline = "#EndOfLine#"
start = "Table bac...'"
kill = "Table batch_blob...'"

#getKeys("db","schadegevallen")
#folder = "DB60324"
#folder = "schaekers_brio"
#read_table("doc_blob_desc.txt",endline)
read_objects("msg_as2_att.txt",endline)
